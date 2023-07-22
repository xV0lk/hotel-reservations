package db

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DBNAME = "hotel-reservations"
	DBURI  = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}

func HandleGetError(c *fiber.Ctx, err error) error {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return err
	}
	return nil
}

func FormatMongoE(e error) string {
	re := regexp.MustCompile(`\{([^}]*)\}`)
	matches := re.FindStringSubmatch(e.Error())
	if len(matches) > 1 {
		return strings.ReplaceAll(strings.TrimSpace(matches[1]), "\"", "'")
	}
	return ""
}
