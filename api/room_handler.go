package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
)

type BookBody struct {
	FromDate  time.Time `json:"fromDate"`
	UntilDate time.Time `json:"untilDate"`
	NumPeople int       `json:"numPeople"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomId := c.Params("id")
	var reqBody BookBody
	err := c.BodyParser(&reqBody)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Printf("-------------------------\nreqBody: %+v\n", reqBody)
	room, err := h.store.Room.GetRoomById(c.Context(), roomId)
	if err != nil {
		fmt.Println("Error: ", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found"})
	}
	fmt.Printf("-------------------------\nroom: %+v\n", room)
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Error"})
	}
	fmt.Printf("-------------------------\nuser: %+v\n", user)
	tPrice := room.BasePrice * float64(reqBody.UntilDate.Sub(reqBody.FromDate).Hours()/24)
	if err = types.ValidateCapacity(reqBody.NumPeople, room); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	cBook := types.Booking{
		UserID:    user.ID,
		RoomID:    room.ID,
		FromDate:  reqBody.FromDate,
		UntilDate: reqBody.UntilDate,
		NumPeople: reqBody.NumPeople,
		Price:     tPrice,
	}
	fmt.Printf("-------------------------\ncBook: %+v\n", cBook)
	return c.Status(fiber.StatusOK).JSON(room)
}
