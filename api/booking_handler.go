package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleMonthBookings(c *fiber.Ctx) error {
	var reqBody types.BookingFilter
	err := c.BodyParser(&reqBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	filter := reqBody.CreateMonthFilter()
	bookings, err := h.store.Booking.FilterBookings(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(bookings)
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.FilterBookings(c.Context(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found"})
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	bookingId := c.Params("id")
	objectId, _ := primitive.ObjectIDFromHex(bookingId)
	bookings, err := h.store.Booking.FilterBookings(c.Context(), bson.M{"_id": objectId})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found"})
	}
	return c.JSON(bookings)
}
