package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/xV0lk/hotel-reservations/api"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbUri    = "mongodb://localhost:27017"
	userColl = "users"
)

var fconfig = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	},
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Load environment variables
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")
	}

	// Initialize handlers
	var (
		// stores
		userStore    = db.NewMongoUserStore(client, db.DBNAME)
		hotelStore   = db.NewMongoHotelStore(client, db.DBNAME)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client, db.DBNAME)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		// handlers
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		// connection
		port = flag.String("port", ":3000", "port to run the server on")
		app  = fiber.New(fconfig)
		// apiV1 = app.Group("/api/v1")
		apiV1 = app.Group("/api/v1", middleware.JWTAuth(userStore))
		admin = apiV1.Group("/admin", middleware.AdminAuth)
	)

	// Create unique email index
	store.User.IndexEmail(context.Background())

	app.Get("/", handleHome)
	// Auth
	app.Post("/api/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	// apiV1.Post("/login", userHandler.HandleLogin)

	// hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/bookings", hotelHandler.HandleGetBookings)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiV1.Get("/hotel/:id/bookings", hotelHandler.HandleGetBookingsById)

	// room handlers
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiV1.Get("/room", roomHandler.HandleGetRooms)

	// booking Handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	apiV1.Get("/booking/month", bookingHandler.HandleMonthBookings)
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Delete("/booking/:id", bookingHandler.HandleCancelBooking)

	app.Listen(*port)
}

// handleHome is a simple handler to test the server
func handleHome(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Server is working!"})
}
