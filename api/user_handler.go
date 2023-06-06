package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/types"
)

// Get all users
func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "John",
		LastName:  "Doe",
	}
	return c.JSON(u)
}

// Get a single user with the id
func HandleGetUser(c *fiber.Ctx) error {
	// We can user fiber.map to create a map[string]interface{}
	return c.JSON(fiber.Map{"user": "This is a single user"})
}
