package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Error"})
	}
	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.Next()
}
