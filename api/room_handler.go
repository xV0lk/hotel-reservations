package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	iutils "github.com/xV0lk/hotel-reservations/utils"
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
	if err != nil {
		return ErrInternal()
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomId := c.Params("id")
	var reqBody types.BookingBody
	err := c.BodyParser(&reqBody)
	if err != nil {
		return ErrBadRequest()
	}
	room, err := h.store.Room.GetRoomById(c.Context(), roomId)
	if err != nil {
		return ErrNotFound()
	}
	if errors := reqBody.Validate(room); len(errors) != 0 {
		return ErrBadRequest()
	}
	// Check if the room is available
	ra, err := h.isRoomAvailable(c.Context(), reqBody, room.ID)
	if err != nil {
		return ErrInternal()
	}
	if !ra {
		return NewError(http.StatusBadRequest, "This room is not available for the time selected")
	}
	// Get user
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrInternal()
	}
	tPrice := room.BasePrice * iutils.Dbd(reqBody.FromDate, reqBody.UntilDate)

	cBook := types.Booking{
		UserID:    user.ID,
		RoomID:    room.ID,
		FromDate:  reqBody.FromDate,
		UntilDate: reqBody.UntilDate,
		NumPeople: reqBody.NumPeople,
		Price:     tPrice,
		Cancelled: false,
	}
	h.store.Booking.InsertBooking(c.Context(), &cBook)
	return c.Status(http.StatusCreated).JSON(cBook)
}

func (h *RoomHandler) isRoomAvailable(ctx *fasthttp.RequestCtx, b types.BookingBody, rId primitive.ObjectID) (bool, error) {
	avFilter := b.CreateAvailabilityFilter(rId)
	cb, err := h.store.Booking.FilterBookings(ctx, avFilter)
	if err != nil {
		return false, ErrInternal()
	}
	available := len(cb) == 0
	return available, nil
}
