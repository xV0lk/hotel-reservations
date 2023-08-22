package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	iutils "github.com/xV0lk/hotel-reservations/utils"
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
	bookings, err := h.store.Booking.GetBookings(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found."})
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	bookingId := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), bookingId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found."})
	}
	if err := bookingAuthorization(c, booking); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err})
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	bookingId := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), bookingId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found"})
	}
	if err := bookingAuthorization(c, booking); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err})
	}
	booking, err = h.store.Booking.CancelBooking(c.Context(), bookingId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.JSON(booking)
}

func bookingAuthorization(c *fiber.Ctx, booking *types.Booking) error {
	user, err := iutils.GetAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID && !user.IsAdmin {
		return fmt.Errorf("Unauthorized")
	}
	return nil
}
