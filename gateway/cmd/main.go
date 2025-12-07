package main

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// jwtCustomClaims are custom claims extending default RegisteredClaims
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func main() {
	e := echo.New()

	// Middleware configurations
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- 1. SETUP PROXY TARGET (Driver Service) ---
	// Get URL from environment variable, default to localhost for local dev
	targetURL := os.Getenv("DRIVER_SERVICE_URL")
	if targetURL == "" {
		targetURL = "http://localhost:8080"
	}
	driverServiceURL, err := url.Parse(targetURL)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Create proxy targets and balancer
	targets := []*middleware.ProxyTarget{{URL: driverServiceURL}}
	balancer := middleware.NewRoundRobinBalancer(targets)

	// --- 2. PUBLIC ROUTES (Accessible by everyone) ---

	// Login endpoint to generate JWT tokens
	e.POST("/login", login)

	// --- 3. PROTECTED ROUTES (Requires valid JWT) ---

	// Group routes starting with /drivers
	r := e.Group("/drivers")

	// Configure JWT middleware
	// "secret" is the signing key. In production, this should come from environment variables.
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}

	// Apply JWT middleware to the group
	r.Use(echojwt.WithConfig(config))

	// If token is valid, forward the request to Driver Service (Reverse Proxy)
	r.Use(middleware.Proxy(balancer))

	// Start Gateway server
	e.Logger.Fatal(e.Start(":8000"))
}

// login handler performs mock authentication and returns a JWT token
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Mock check (In a real app, check against database)
	// Expected credentials -> username: admin, password: password123
	if username != "admin" || password != "password123" {
		return echo.ErrUnauthorized
	}

	// Set custom claims for the token
	claims := &jwtCustomClaims{
		"Bitaksi Admin",
		true,
		jwt.RegisteredClaims{
			// Token expires in 72 hours
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims using HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token using the secret key
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	// Return token in JSON response
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
