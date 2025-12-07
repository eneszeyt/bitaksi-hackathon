package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eneszeyt/bitaksi-driver-service/internal/config"
)

func main() {

	// 1. load settings

	cfg := config.LoadConfig()

	fmt.Printf("Driver Service is initializing...\nPORT: %s\nDB: %s\n", cfg.Port, cfg.DBName)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Driver Service (Structured) is Working! 🚕")
	})

	// initialize to server

	addr := ":" + cfg.Port
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server Error : %v", err)
	}
}
