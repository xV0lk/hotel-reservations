package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var body AuthParams
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	user, err := h.userStore.GetUser(c.Context(), bson.M{"email": body.Email})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found"})
	}
	if !types.IsValidPassword(body.Password, user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Password"})
	}
	token, err := createUserToken(user)
	if err != nil {
		fmt.Println("Error: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal error"})
	}
	resp := AuthResponse{
		User:  user,
		Token: token,
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func createUserToken(user *types.User) (string, error) {
	now := time.Now()
	expire := now.Add(time.Hour * 24).Unix()
	claims := jwt.MapClaims{
		"id":         user.ID,
		"email":      user.Email,
		"expiration": expire,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	ts, err := token.SignedString([]byte(secret))
	return ts, err
}
