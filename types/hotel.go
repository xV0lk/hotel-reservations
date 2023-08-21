package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   float64              `bson:"rating" json:"rating"`
}

type HotelBookings struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   float64              `bson:"rating" json:"rating"`
	Bookings []Booking            `bson:"bookings" json:"bookings"`
}

type RoomType int

const (
	_ RoomType = iota
	Single
	Double
	SeaSide
	Deluxe
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type,omitempty" json:"type,omitempty"`
	BasePrice float64            `bson:"basePrice,omitempty" json:"basePrice,omitempty"`
	HotelId   primitive.ObjectID `bson:"hotelId,omitempty" json:"hotelId,omitempty"`
}
