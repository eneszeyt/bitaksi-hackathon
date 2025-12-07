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

	// middleware configurations
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// cors middleware configuration to allow frontend access
	// this is crucial for the react app to communicate with the backend
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // allow all origins (restrict this in production)
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// --- 1. setup proxy target (driver service) ---
	// get url from environment variable, default to localhost for local dev
	targetURL := getEnv("DRIVER_SERVICE_URL", "http://localhost:8080")

	driverServiceURL, err := url.Parse(targetURL)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// create proxy targets and balancer
	targets := []*middleware.ProxyTarget{{URL: driverServiceURL}}
	balancer := middleware.NewRoundRobinBalancer(targets)

	// --- 2. public routes (accessible by everyone) ---

	// login endpoint to generate jwt tokens
	e.POST("/login", login)

	// --- 3. protected routes (requires valid jwt) ---

	// group routes starting with /drivers
	r := e.Group("/drivers")

	// security fix: read secret key from environment variable
	jwtSecret := getEnv("JWT_SECRET", "secret")

	// configure jwt middleware
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(jwtSecret),
	}

	// apply jwt middleware to the group
	r.Use(echojwt.WithConfig(config))

	// if token is valid, forward the request to driver service (reverse proxy)
	r.Use(middleware.Proxy(balancer))

	// start gateway server
	e.Logger.Fatal(e.Start(":8000"))
}

// login handler performs mock authentication and returns a jwt token
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// mock check (in a real app, check against database)
	// expected credentials -> username: admin, password: password123
	if username != "admin" || password != "password123" {
		return echo.ErrUnauthorized
	}

	// set custom claims for the token
	claims := &jwtCustomClaims{
		"Bitaksi Admin",
		true,
		jwt.RegisteredClaims{
			// token expires in 72 hours
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// create token with claims using hs256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// security fix: sign token using the same secret from environment
	jwtSecret := getEnv("JWT_SECRET", "secret")

	// generate encoded token
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return err
	}

	// return token in json response
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

// getEnv retrieves the value of the environment variable named by the key.
// it returns the value, which will be empty if the variable is not present.
// if the variable is not present, it returns the fallback value.

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
