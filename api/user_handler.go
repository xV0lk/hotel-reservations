package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

// Get all users
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var id = c.Params("id")
	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// Get a single user with the id
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// We can user fiber.map to create a map[string]interface{}
	return c.JSON(users)
}
