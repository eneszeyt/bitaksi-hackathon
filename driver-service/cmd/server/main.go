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

	// connect to database
	mongoClient, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connection failed: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// dependency injection
	db := mongoClient.Database(cfg.DBName)
	repo := repository.NewDriverRepository(db)
	svc := service.NewDriverService(repo)
	h := handler.NewDriverHandler(svc)

	// --- ROUTES ---

	// 1. /drivers -> GET (List) & POST (Create)
	http.HandleFunc("/drivers", h.DriversRoot)

	// 2. /drivers/nearby -> GET (Nearby Search)
	// IMPORTANT: This must be defined BEFORE "/drivers/" because "/drivers/" matches all subpaths
	http.HandleFunc("/drivers/nearby", h.SearchNearby)

	// 3. /drivers/ -> PUT (Update) (trailing slash matches subpaths like /drivers/123)
	http.HandleFunc("/drivers/", h.DriverByID)

	// start server
	addr := ":" + cfg.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
