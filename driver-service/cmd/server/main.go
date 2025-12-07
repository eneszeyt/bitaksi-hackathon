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

	_ "github.com/eneszeyt/bitaksi-driver-service/docs" // This line is crucial for swagger to find generated docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Bitaksi Driver Service API
// @version         1.0
// @description     This is a sample driver service for Bitaksi Hackathon.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

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

	// 1. Swagger Documentation Route
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// 2. /drivers -> GET (List) & POST (Create)
	http.HandleFunc("/drivers", h.DriversRoot)

	// 3. /drivers/nearby -> GET (Nearby Search)
	http.HandleFunc("/drivers/nearby", h.SearchNearby)

	// 4. /drivers/ -> PUT (Update)
	http.HandleFunc("/drivers/", h.DriverByID)

	// start server
	addr := ":" + cfg.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
