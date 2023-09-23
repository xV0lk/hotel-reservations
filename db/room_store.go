package db

import (
	"context"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) error
	InsertManyRooms(ctx context.Context, rooms []types.Room, hId primitive.ObjectID) error
	GetRooms(ctx *fasthttp.RequestCtx, filter bson.M) ([]*types.Room, error)
	GetRoomById(ctx *fasthttp.RequestCtx, id string) (*types.Room, error)
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

func (s *MongoRoomStore) InsertManyRooms(ctx context.Context, rooms []types.Room, hId primitive.ObjectID) error {
	var documents []interface{}
	for _, room := range rooms {
		room.HotelId = hId
		documents = append(documents, room)
	}
	iRooms, err := s.coll.InsertMany(ctx, documents)
	if err != nil {
		return err
	}
	for _, iRoom := range iRooms.InsertedIDs {
		filter := bson.M{"_id": hId}
		update := bson.M{"$push": bson.M{"rooms": iRoom.(primitive.ObjectID)}}
		if _, err = s.HotelStore.Update(ctx, filter, update); err != nil {
			log.Fatal(err)
		}
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

func (s *MongoRoomStore) GetRoomById(ctx *fasthttp.RequestCtx, id string) (*types.Room, error) {
	var room types.Room
	objectId, _ := primitive.ObjectIDFromHex(id)
	if err := s.coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&room); err != nil {
		return nil, err
	}
	return &room, nil
}
