package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID    primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	FromDate  time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	UntilDate time.Time          `bson:"untilDate,omitempty" json:"untilDate,omitempty"`
	Price     float64            `bson:"price,omitempty" json:"price,omitempty"`
	NumPeople int                `bson:"numPeople,omitempty" json:"numPeople,omitempty"`
}

func ValidateCapacity(p int, r *Room) error {
	switch r.Type {
	case Single:
		if p > 1 {
			return fmt.Errorf("single room can only accommodate 1 person, but got %d", p)
		}
	case Double:
		if p > 3 {
			return fmt.Errorf("double room can only accommodate 3 people, but got %d", p)
		}
	case SeaSide:
		if p > 3 {
			return fmt.Errorf("sea-side room can only accommodate 3 people, but got %d", p)
		}
	case Deluxe:
		if p > 4 {
			return fmt.Errorf("deluxe room can only accommodate 4 people, but got %d", p)
		}
	default:
		return fmt.Errorf("unknown room type: %d", r.Type)
	}
	return nil
}
