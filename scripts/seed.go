package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore *db.MongoHotelStore
	roomStore  *db.MongoRoomStore
	ctx        = context.Background()
)

func seedHotel(name, location string, rating, priceCategory float64) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}

	rooms := []types.Room{
		{
			Type:      types.Single,
			BasePrice: (priceCategory * 50) - 1,
		},
		{
			Type:      types.Deluxe,
			BasePrice: (priceCategory * 80) - 1,
		},
		{
			Type:      types.Double,
			BasePrice: (priceCategory * 110) - 1,
		},
		{
			Type:      types.SeaSide,
			BasePrice: (priceCategory * 140) - 1,
		},
	}
	if err := hotelStore.InsertHotel(ctx, &hotel); err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelId = hotel.ID
		if err := roomStore.InsertRoom(ctx, &room); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("The Coffin", "Transylvania", 3.5, 1)
	seedHotel("Yokai Inn", "Japan", 4.2, 2)
	seedHotel("Sherlock hideout", "London", 4.9, 2.5)
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to seed MongoDB!")

	// we need to drop the database first
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
