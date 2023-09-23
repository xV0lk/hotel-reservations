package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	iutils "github.com/xV0lk/hotel-reservations/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, n, l string, a bool) *types.User {
	newUser := types.NewUserParams{
		FirstName: n,
		LastName:  l,
		Email:     fmt.Sprintf("%s@%s.com", n, l),
		Password:  fmt.Sprintf("%s_%s_P4$$", n, l),
		IsAdmin:   a,
	}
	if errors := newUser.Validate(); len(errors) != 0 {
		log.Fatal(errors)
	}
	user, err := types.NewUserFromParams(&newUser)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.User.InsertUser(context.TODO(), user); err != nil {
		log.Fatal(err)
	}
	return user
}

func AddRoom(store *db.Store, t types.RoomType, p int, h primitive.ObjectID) *types.Room {
	ctx := context.Background()
	room := &types.Room{
		Type:      t,
		BasePrice: p,
		HotelId:   h,
	}
	if err := store.Room.InsertRoom(ctx, room); err != nil {
		log.Fatal(err)
	}
	return room
}

func AddERoom(store *db.Store, r types.Room) {
	ctx := context.Background()
	if err := store.Room.InsertRoom(ctx, &r); err != nil {
		log.Fatal(err)
	}
}

func AddHotel(store *db.Store, name, location string, rating float64, priceCategory int) *types.Hotel {
	ctx := context.Background()
	hotel := &types.Hotel{
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
	if err := store.Hotel.InsertHotel(ctx, hotel); err != nil {
		log.Fatal(err)
	}
	if err := store.Room.InsertManyRooms(ctx, rooms, hotel.ID); err != nil {
		log.Fatal(err)
	}
	hotel, _ = store.Hotel.GetHotelById(ctx, hotel.ID.Hex())
	return hotel
}

func AddBooking(store *db.Store, user primitive.ObjectID, room *types.Room, from, till time.Time, guests int) *types.Booking {
	ctx := context.Background()
	booking := &types.Booking{
		UserID:    user,
		RoomID:    room.ID,
		FromDate:  from,
		UntilDate: till,
		NumPeople: guests,
		Price:     room.BasePrice * iutils.Dbd(from, till),
	}
	if err := store.Booking.InsertBooking(ctx, booking); err != nil {
		log.Fatal(err)
	}
	return booking
}
