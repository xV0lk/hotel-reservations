package db

import (
	"context"

	"github.com/valyala/fasthttp"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) error
	GetRooms(ctx *fasthttp.RequestCtx, filter bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) error {
	result, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return err
	}
	room.ID = result.InsertedID.(primitive.ObjectID)

	// add room to hotel
	filter := bson.M{"_id": room.HotelId}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	_, err = s.HotelStore.Update(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (s *MongoRoomStore) GetRooms(ctx *fasthttp.RequestCtx, filter bson.M) ([]*types.Room, error) {
	var rooms []*types.Room
	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
