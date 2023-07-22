package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	var id = c.Params("id")
	hotel, err := h.hotelStore.GetHotelById(c.Context(), id)
	db.HandleGetError(c, err)
	return c.JSON(hotel)
}

type HotelParams struct {
	Name   string
	Rating float64
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelParams
	if err := c.QueryParser(&params); err != nil {
		return err
	}
	fmt.Printf("-------------------------\nhotelParams: %+v\n", params)
	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(hotels)
}
