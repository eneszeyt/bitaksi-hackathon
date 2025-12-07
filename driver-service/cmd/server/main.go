package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/eneszeyt/bitaksi-driver-service/internal/config"
	"github.com/eneszeyt/bitaksi-driver-service/internal/handler"
	"github.com/eneszeyt/bitaksi-driver-service/internal/repository"
	"github.com/eneszeyt/bitaksi-driver-service/internal/service"
	"github.com/eneszeyt/bitaksi-driver-service/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("starting driver service on port %s...\n", cfg.Port)

	mongoClient, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connection failed: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(cfg.DBName)
	repo := repository.NewDriverRepository(db)
	svc := service.NewDriverService(repo)
	h := handler.NewDriverHandler(svc)

	// --- ROUTES ---
	// 1. Create (POST /drivers)
	http.HandleFunc("/drivers", h.CreateDriver)

	// 2. Update (PUT /drivers/{id})
	// Sondaki "/" işareti önemli, bu sayede /drivers/123 gibi alt yolları yakalar
	http.HandleFunc("/drivers/", h.UpdateDriver)

	addr := ":" + cfg.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
