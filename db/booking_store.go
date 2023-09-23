package db

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/xV0lk/hotel-reservations/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) error
	FilterBookings(ctx *fasthttp.RequestCtx, filter bson.M) ([]*types.Booking, error)
	GetBookings(ctx *fasthttp.RequestCtx) ([]*types.Booking, error)
	GetBookingById(ctx *fasthttp.RequestCtx, id string) (*types.Booking, error)
	CancelBooking(ctx *fasthttp.RequestCtx, id string) (*types.Booking, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client, dbname string) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) error {
	result, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return err
	}
	booking.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *MongoBookingStore) GetBookings(ctx *fasthttp.RequestCtx) ([]*types.Booking, error) {
	var bookings []*types.Booking
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) FilterBookings(ctx *fasthttp.RequestCtx, filter bson.M) ([]*types.Booking, error) {
	var bookings []*types.Booking
	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) GetBookingById(ctx *fasthttp.RequestCtx, id string) (*types.Booking, error) {
	var booking *types.Booking
	oId, _ := primitive.ObjectIDFromHex(id)
	if err := s.coll.FindOne(ctx, bson.M{"_id": oId}).Decode(&booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *MongoBookingStore) CancelBooking(ctx *fasthttp.RequestCtx, id string) (*types.Booking, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"cancelled": true}}
	result, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no booking with id %s", id)
	}
	if updated, err := s.GetBookingById(ctx, id); err != nil {
		return nil, err
	} else {
		return updated, nil
	}
}
