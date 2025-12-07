package main

import (
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Create Echo instance
	e := echo.New()

	// 2. Middleware
	e.Use(middleware.Logger())  // Log requests
	e.Use(middleware.Recover()) // Recover from panics

	// 3. Setup Proxy Target
	// Get URL from environment variable (for Docker), default to localhost (for local dev)
	targetURL := os.Getenv("DRIVER_SERVICE_URL")
	if targetURL == "" {
		targetURL = "http://localhost:8080"
	}

	driverServiceURL, err := url.Parse(targetURL)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// 4. Configure Proxy
	targets := []*middleware.ProxyTarget{
		{
			URL: driverServiceURL,
		},
	}

	// 5. Define Routes
	// Group /drivers routes and forward them to the driver service
	e.Group("/drivers", middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	// 6. Start Gateway on port 8000
	e.Logger.Fatal(e.Start(":8000"))
}
