package types

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID    primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	FromDate  time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	UntilDate time.Time          `bson:"untilDate,omitempty" json:"untilDate,omitempty"`
	Price     int                `bson:"price,omitempty" json:"price,omitempty"`
	NumPeople int                `bson:"numPeople,omitempty" json:"numPeople,omitempty"`
	Cancelled bool               `bson:"cancelled" json:"cancelled"`
}

type BookingBody struct {
	FromDate  time.Time `json:"fromDate"`
	UntilDate time.Time `json:"untilDate"`
	NumPeople int       `json:"numPeople"`
}

type BookingFilter struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

func (b BookingBody) Validate(r *Room) map[string]string {
	errors := map[string]string{}
	// capacity validation
	if err := validateCapacity(b.NumPeople, r); err != nil {
		errors["capacity"] = err.Error()
	}
	// Date validation
	if err := validateDate(b); len(err) != 0 {
		errors["date"] = strings.Join(err, ", ")
	}
	// Availability Validation

	return errors
}

func validateDate(b BookingBody) []string {
	errors := []string{}
	now := time.Now()
	if now.After(b.FromDate) {
		errors = append(errors, "can't use a date before today date as starting date")
	}
	if now.After(b.UntilDate) {
		errors = append(errors, "can't use a date before today as ending date")
	}
	if b.FromDate.After(b.UntilDate) {
		errors = append(errors, "end date must be after start date")
	}
	return errors
}

func validateCapacity(p int, r *Room) error {
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

func (b BookingFilter) CreateMonthFilter() bson.M {
	fromDateConditions := bson.M{
		"$and": []bson.M{
			{"$eq": bson.A{bson.M{"$month": "$fromDate"}, b.Month}},
			{"$eq": bson.A{bson.M{"$year": "$fromDate"}, b.Year}},
		},
	}

	untilDateConditions := bson.M{
		"$and": []bson.M{
			{"$eq": bson.A{bson.M{"$month": "$untilDate"}, b.Month}},
			{"$eq": bson.A{bson.M{"$year": "$untilDate"}, b.Year}},
		},
	}

	return bson.M{
		"cancelled": false,
		"$expr": bson.M{
			"$or": bson.A{fromDateConditions, untilDateConditions},
		},
	}
}

func (b BookingBody) CreateAvailabilityFilter(rId primitive.ObjectID) bson.M {
	return bson.M{
		"roomID":    rId,
		"cancelled": false,
		"$or": []bson.M{
			{"fromDate": bson.M{"$gte": b.FromDate, "$lt": b.UntilDate}},
			{"untilDate": bson.M{"$gt": b.FromDate, "$lte": b.UntilDate}},
			{"fromDate": bson.M{"$lte": b.FromDate}, "untilDate": bson.M{"$gte": b.UntilDate}},
		},
	}
}
