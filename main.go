package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

// var (
// 	appid  = "wx16bb3389273e4759"
// 	secret = "2e6b70686ce2e0fa9f41b31677acc8a2"
// )

var (
	appid  = "wx370126c8bcf8d00c"
	secret = "e2bd2db2f82b824c66d021b6d4f5b7bb"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Post("/api/pc/import", ExcelImport)

	app.Post("/api/mini/register", Register)

	app.Get("/api/mini/getUserInfo", GetUserInfo)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
