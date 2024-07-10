package main

import (
	"bleve-proj/internal/repositories"
	"bleve-proj/internal/routes"
	"log"
	"os"

	"bleve-proj/cmd/server/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize index directory
	indexDir := os.Getenv("INDEX_DIRECTORY")
	if indexDir == "" {
		log.Fatal("INDEX_DIRECTORY environment variable is not set")
	}

	// Initialize index repository with index directory
	repositories.IndexRepo = repositories.NewIndexRepository(indexDir)

	// Create a new Fiber instance
	app := fiber.New()

	prefixPath := os.Getenv("PREFIX_PATH")
	basePath := "/api/v1"
	swaggerUrl := "/swagger/doc.json"
	if prefixPath != "" {
		basePath = "/" + prefixPath + basePath
		swaggerUrl = "/" + prefixPath + swaggerUrl
	}
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: swaggerUrl,
	}))

	// Define routes
	routes.InitRoutes(app)

	// Start the server
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
	docs.SwaggerInfo.BasePath = basePath
	log.Fatal(app.Listen(":" + port))
}
