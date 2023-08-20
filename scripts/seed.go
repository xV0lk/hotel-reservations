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
	userStore  *db.MongoUserStore
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
			Type:      types.Double,
			BasePrice: (priceCategory * 80) - 1,
		},
		{
			Type:      types.SeaSide,
			BasePrice: (priceCategory * 110) - 1,
		},
		{
			Type:      types.Deluxe,
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

func seedUser(name, lName, email, password string) {
	newUser := types.NewUserParams{
		FirstName: name,
		LastName:  lName,
		Email:     email,
		Password:  password,
	}
	if errors := newUser.Validate(); len(errors) != 0 {
		log.Fatal(errors)
	}
	user, err := types.NewUserFromParams(&newUser)
	if err != nil {
		log.Fatal(err)
	}
	if err := userStore.InsertUser(ctx, user); err != nil {
		log.Fatal(err)
	}
}

func main() {
	seedHotel("The Coffin", "Transylvania", 3.5, 1)
	seedHotel("Yokai Inn", "Japan", 4.2, 2)
	seedHotel("Sherlock hideout", "London", 4.9, 2.5)
	seedUser("Jorge", "Rojas", "jorge.otto.415@gmail.com", "Jrojas1234$")
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
	userStore = db.NewMongoUserStore(client, db.DBNAME)
}
