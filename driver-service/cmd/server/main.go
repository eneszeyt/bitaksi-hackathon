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
	fmt.Printf("Driver Service Başlıyor... Port: %s\n", cfg.Port)

	// 1. Veritabanına Bağlan
	mongoClient, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("❌ Veritabanına bağlanılamadı: %v", err)
	}

	// Uygulama kapanırken bağlantıyı kes
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Disconnect hatası: %v", err)
		}
	}()

	// 2. HTTP Sunucusunu Başlat
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Driver Service + MongoDB Bağlantısı Aktif! 🚀")
	})

	addr := ":" + cfg.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Sunucu hatası: %v", err)
	}
}
