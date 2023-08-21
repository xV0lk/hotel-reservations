package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	db.HandleGetError(c, err)
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomId := c.Params("id")
	var reqBody types.BookingBody
	err := c.BodyParser(&reqBody)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	room, err := h.store.Room.GetRoomById(c.Context(), roomId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found"})
	}
	if errors := reqBody.Validate(room); len(errors) != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors})
	}
	// Check if the room is available
	ra, err := h.isRoomAvailable(c.Context(), reqBody, room.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	if !ra {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "This room is not available for the time selected"})
	}
	// Get user
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Error"})
	}
	tPrice := room.BasePrice * float64(reqBody.UntilDate.Sub(reqBody.FromDate).Hours()/24)

	cBook := types.Booking{
		UserID:    user.ID,
		RoomID:    room.ID,
		FromDate:  reqBody.FromDate,
		UntilDate: reqBody.UntilDate,
		NumPeople: reqBody.NumPeople,
		Price:     tPrice,
	}
	h.store.Booking.InsertBooking(c.Context(), &cBook)
	return c.Status(fiber.StatusCreated).JSON(cBook)
}

func (h *RoomHandler) isRoomAvailable(ctx *fasthttp.RequestCtx, b types.BookingBody, rId primitive.ObjectID) (bool, error) {
	avFilter := b.CreateAvailabilityFilter(rId)
	cb, err := h.store.Booking.FilterBookings(ctx, avFilter)
	if err != nil {
		return false, err
	}
	available := len(cb) == 0
	return available, nil
}
