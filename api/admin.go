package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrInternal()
	}
	if !user.IsAdmin {
		return ErrForbidden()
	}
	return c.Next()
}
