package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/mongo"
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
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
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

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var id = c.Params("id")
	err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully!"})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var newUser *types.NewUserParams
	if err := c.BodyParser(&newUser); err != nil {
		return err
	}
	if errors := newUser.Validate(); len(errors) != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors})
	}
	user, err := types.NewUserFromParams(newUser)
	if err != nil {
		return err
	}
	err = h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		id       = c.Params("id")
		jsonData map[string]any
		newUser  *types.NewUserParams
	)
	if err := c.BodyParser(&newUser); err != nil {
		return err
	}
	if errors := newUser.Validate(); len(errors) != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	err := c.BodyParser(&jsonData)
	if err != nil {
		return err
	}
	if err := types.CheckUserBody(jsonData); err != nil {
		return err
	}
	updated, err := h.userStore.UpdateUser(c.Context(), id, jsonData)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(updated)
}
