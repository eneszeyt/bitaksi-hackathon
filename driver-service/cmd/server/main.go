package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/eneszeyt/bitaksi-driver-service/internal/config"
	"github.com/eneszeyt/bitaksi-driver-service/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("Driver Service is starting... Port: %s\n", cfg.Port)

	// 1. connect to database

	mongoClient, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// disconnect when app closes

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Disconnect error: %v", err)
		}
	}()

	// 2. initialize to HTTP Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Driver Service + MongoDB connection is Active! 🚀")
	})

	addr := ":" + cfg.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
