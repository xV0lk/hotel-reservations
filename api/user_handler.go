package api

import (
	"net/http"

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
		return ErrNotFound()
	}
	return c.JSON(user)
}

// Get a single user with the id
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrInternal()
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var id = c.Params("id")
	err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		return ErrBadRequest()
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully!"})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var newUser *types.NewUserParams
	if err := c.BodyParser(&newUser); err != nil {
		return err
	}
	if errors := newUser.Validate(); len(errors) != 0 {
		return NewMapError(http.StatusBadRequest, errors)
	}
	if err := validateAdminCreation(c, newUser); err != nil {
		return NewError(http.StatusTeapot, err.Error())
	}
	user, err := types.NewUserFromParams(newUser)
	if err != nil {
		return ErrInternal()
	}
	err = h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return NewError(http.StatusBadRequest, "Email already exists")
		}
		return ErrBadRequest()
	}
	return c.Status(http.StatusCreated).JSON(user)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		id         = c.Params("id")
		jsonData   map[string]any
		updateUser *types.UpdateUserParams
	)
	if err := c.BodyParser(&updateUser); err != nil {
		return ErrBadRequest()
	}
	err := c.BodyParser(&jsonData)
	if err != nil {
		return ErrBadRequest()
	}
	if err := updateUser.CheckBody(jsonData); err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}
	if errors := updateUser.Validate(); len(errors) != 0 {
		return NewMapError(http.StatusBadRequest, errors)
	}
	updated, err := h.userStore.UpdateUser(c.Context(), id, updateUser)
	if err != nil {
		return ErrInternal()
	}
	return c.JSON(updated)
}

func validateAdminCreation(c *fiber.Ctx, nu *types.NewUserParams) error {
	// Get user
	iUser, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrInternal()
	}
	if !iUser.IsAdmin && nu.IsAdmin {
		return ErrForbidden()
	}
	return nil
}
