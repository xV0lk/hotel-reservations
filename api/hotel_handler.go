package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	var id = c.Params("id")
	hotel, err := h.store.Hotel.GetHotelById(c.Context(), id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	var id = c.Params("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{"hotelId": objectId})
	db.HandleGetError(c, err)
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetBookingsById(c *fiber.Ctx) error {
	var id = c.Params("id")
	bookings, err := h.store.Hotel.GetHotelBookings(c.Context(), id)
	db.HandleGetError(c, err)
	return c.Status(fiber.StatusOK).JSON(bookings)
}

func (h *HotelHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Hotel.GetHotelBookings(c.Context(), "")
	db.HandleGetError(c, err)
	return c.Status(fiber.StatusOK).JSON(bookings)
}
