package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bleve-proj/cmd/restapi/server/docs"
	"bleve-proj/internal/repositories"
	"bleve-proj/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

var BasePath = "api/v1"

func main() {
	// Initialize environment variables
	initEnv()

	// Initialize repositories
	initRepositories()

	// Create a new Fiber instance and initialize routes
	app := initFiberApp()

	// Start the server
	startServer(app)

	// Call the graceful shutdown function
	gracefulShutdown(app)
}

// initEnv loads environment variables from .env file
func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// initRepositories initializes the index repository
func initRepositories() {
	indexDir := os.Getenv("INDEX_DIRECTORY")
	if indexDir == "" {
		log.Fatal("INDEX_DIRECTORY environment variable is not set")
	}

	// Initialize index repository with index directory
	repositories.IndexRepo = repositories.NewIndexRepository(indexDir)
}

// initFiberApp creates a new Fiber instance and sets up routes
func initFiberApp() *fiber.App {
	app := fiber.New()

	prefixPath := os.Getenv("PREFIX_PATH")
	swaggerUrl := "/swagger/doc.json"
	if prefixPath != "" {
		BasePath = "/" + prefixPath + BasePath
		swaggerUrl = "/" + prefixPath + swaggerUrl
	}
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: swaggerUrl,
	}))

	// Define routes
	routes.InitRoutes(app)

	return app
}

// startServer starts the Fiber server in a separate goroutine
func startServer(app *fiber.App) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	domain := os.Getenv("DOMAIN")
	host := "localhost:" + port
	if domain != "" {
		host = domain
	}
	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.Title = "BLEVE"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = BasePath

	// Run server in a separate goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Fiber server failed to start: %v", err)
		}
	}()
}

// gracefulShutdown handles the graceful shutdown of the server
func gracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Server().ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close the index connections
	repositories.IndexRepo.Close()

	log.Println("Server exiting")
}
