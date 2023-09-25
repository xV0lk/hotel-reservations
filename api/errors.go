package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		if apiError.Map != nil {
			return c.Status(apiError.Code).JSON(fiber.Map{"error": apiError.Map})
		}
		return c.Status(apiError.Code).JSON(fiber.Map{"error": apiError.Err})
	}
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(fiber.Map{"error": apiError.Err})
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
}

type Error struct {
	Code int               `json:"code"`
	Err  string            `json:"error"`
	Map  map[string]string `json:"map"`
}

// Error implements the error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func NewMapError(code int, m map[string]string) Error {
	return Error{
		Code: code,
		Map:  m,
	}
}

func ErrNotFound() Error {
	return NewError(http.StatusNotFound, "The id you provided is invalid")
}

func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "Bad Request")
}

func ErrInternal() Error {
	return NewError(http.StatusInternalServerError, "Internal Server Error")
}

func ErrUnauthorized() Error {
	return NewError(http.StatusUnauthorized, "Unauthorized")
}

func ErrForbidden() Error {
	return NewError(http.StatusForbidden, "You don't have permission to access this resource")
}
