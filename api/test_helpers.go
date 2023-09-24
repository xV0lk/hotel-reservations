package api

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDbUri  = "mongodb://localhost:27017"
	testDbName = "hotel-reservations-test"
)

type expected struct {
	status int
	body   any
}

type postTest[T any] struct {
	name  string
	ttype string
	input T
	expected
}

type failUserResponse struct {
	Error map[string]string `json:"error"`
}

type userTest types.NewUserParams

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) Drop(t *testing.T) {
	if err := tdb.client.Database(testDbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbUri))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to test MongoDB!")
	hotelStore := db.NewMongoHotelStore(client, testDbName)
	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client, testDbName),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client, testDbName),
		},
	}
}
