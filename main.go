package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/api"
	"github.com/xV0lk/hotel-reservations/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbUri    = "mongodb://localhost:27017"
	userColl = "users"
)

var fconfig = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	},
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Initialize handlers
	var (
		// stores
		userStore  = db.NewMongoUserStore(client, db.DBNAME)
		hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		store      = &db.Store{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}
		// handlers
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		// connection
		port  = flag.String("port", ":3000", "port to run the server on")
		app   = fiber.New(fconfig)
		apiV1 = app.Group("/api/v1")
	)

	app.Get("/", handleHome)

	// user handlers
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*port)
}

// handleHome is a simple handler to test the server
func handleHome(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Server is working!"})
}
