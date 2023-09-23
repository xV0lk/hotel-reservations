package iutils

import (
	"fmt"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/types"
)

func GetAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}
	return user, nil
}

func ValidateAdmin(c *fiber.Ctx) error {
	user, err := GetAuthUser(c)
	if err != nil {
		return fmt.Errorf("Unauthorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("Unauthorized")
	}
	return nil
}

func Dbd(from, till time.Time) int {
	return int(math.Round(till.Sub(from).Hours() / 24))
}
