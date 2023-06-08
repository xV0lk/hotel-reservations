package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
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

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var newUser *types.User
	if err := c.BodyParser(&newUser); err != nil {
		return err
	}
	user, err := h.userStore.PostUser(c.Context(), newUser)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
