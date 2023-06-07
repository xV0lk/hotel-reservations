package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/api"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbUri    = "mongodb://localhost:27017"
	DBNAME   = "hotel-reservations"
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
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	ctx := context.Background()
	coll := client.Database(DBNAME).Collection(userColl)

	user := types.User{
		FirstName: "Jorge",
		LastName:  "Rojas",
	}

	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	var resUser types.User
	if err := coll.FindOne(ctx, bson.M{}).Decode(&resUser); err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println("resUser: ", resUser)

	port := flag.String("port", ":3000", "port to run the server on")
	app := fiber.New(fconfig)
	app.Get("/", handleHome)

	// we can create api groups
	apiV1 := app.Group("/api/v1")
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)

	app.Listen(*port)
}

// handleHome is a simple handler to test the server
func handleHome(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "Server is working!"})
}
