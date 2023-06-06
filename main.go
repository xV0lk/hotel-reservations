package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/api"
)

func main() {
	port := flag.String("port", ":3000", "port to run the server on")
	app := fiber.New()
	app.Get("/", handleHome)

	// we can create api groups
	apiV1 := app.Group("/api/v1")
	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*port)
}

// handleHome is a simple handler to test the server
func handleHome(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Server is working!"})
}
