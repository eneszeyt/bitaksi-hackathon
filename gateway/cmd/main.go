package main

import (
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Create Echo instance (Corrected: lowercase 'echo')
	e := echo.New()

	// 2. Middleware
	e.Use(middleware.Logger())  // Log requests
	e.Use(middleware.Recover()) // Recover from panics

	// 3. Setup Proxy Target (Driver Service)
	// We want to forward requests to localhost:8080
	driverServiceURL, err := url.Parse("http://localhost:8080")
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
	// balancer.RoundRobin is a load balancing strategy
	e.Group("/drivers", middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	// 6. Start Gateway on port 8000
	e.Logger.Fatal(e.Start(":8000"))
}
