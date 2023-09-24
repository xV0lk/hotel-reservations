package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/xV0lk/hotel-reservations/api"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/db/fixtures"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to seed MongoDB!")

	// we need to drop the database first
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	store := &db.Store{
		Hotel:   hotelStore,
		User:    db.NewMongoUserStore(client, db.DBNAME),
		Booking: db.NewMongoBookingStore(client, db.DBNAME),
		Room:    db.NewMongoRoomStore(client, hotelStore),
	}
	admin := fixtures.AddUser(store, "Jorge", "Rojas", true)
	adminToken, _ := api.CreateUserToken(admin)
	fmt.Printf("-------------------------\nadmin: %s\n", adminToken)
	user := fixtures.AddUser(store, "John", "Doe", false)
	userToken, _ := api.CreateUserToken(user)
	fmt.Printf("-------------------------\nuser: %s\n", userToken)
	hotel := fixtures.AddHotel(store, "The Coffin", "Transylvania", 3.5, 1)
	room := fixtures.AddRoom(store, types.Double, 155, hotel.ID)
	booking := fixtures.AddBooking(store, admin.ID, room, time.Now(), time.Now().AddDate(0, 0, 5), 2)
	fixtures.AddBooking(store, user.ID, room, time.Now().AddDate(0, 0, 6), time.Now().AddDate(0, 0, 7), 1)
	bookingH, _ := json.MarshalIndent(booking, "", "  ")
	fmt.Printf("-------------------------\nbooking: %s\n", string(bookingH))
	fixtures.AddHotel(store, "Yokai Inn", "Japan", 4.2, 2)
	fixtures.AddHotel(store, "Sherlock hideout", "London", 4.9, 3)
}
