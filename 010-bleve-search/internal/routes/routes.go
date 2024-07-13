package routes

import (
	"bleve-proj/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {

	app.Get("/api/v1/", controllers.GetPartitions)
	app.Post("/api/v1/:partition/create", controllers.CreateIndex)
	app.Post("/api/v1/:partition/index", controllers.IndexDocument)
	app.Post("/api/v1/:partition/bulk", controllers.BulkLoadDocuments)
	app.Get("/api/v1/:partition/search", controllers.SearchDocument)
	app.Get("/api/v1/:partition/count", controllers.GetDocumentCount)
	app.Get("/api/v1/:partition/check", controllers.CheckIndexContents)
	app.Get("/api/v1/:partition/dictionary", controllers.PrintTermDictionary)
	app.Get("/api/v1/:partition/dump", controllers.DumpIndexContents)
	app.Get("/api/v1/:partition/fields", controllers.ListIndexFields)
	app.Get("/api/v1/:partition/mapping", controllers.PrintIndexMapping)
	app.Delete("/api/v1/:partition/delete", controllers.DeleteIndex)
	// Add more routes as needed
}
