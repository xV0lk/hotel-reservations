package db

import (
	"context"

	"github.com/valyala/fasthttp"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	InsertHotel(ctx context.Context, hotel *types.Hotel) error
	Update(ctx context.Context, filter, update bson.M) (*types.Hotel, error)
	GetHotelById(ctx context.Context, id string) (*types.Hotel, error)
	GetHotels(ctx *fasthttp.RequestCtx) ([]*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) error {
	result, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return err
	}
	hotel.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter, update bson.M) (*types.Hotel, error) {
	result, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	hotelId := filter["_id"].(primitive.ObjectID).Hex()
	hotel, err := s.GetHotelById(ctx, hotelId)
	if err != nil {
		return nil, err
	}
	return hotel, nil
}

func (s *MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	var hotel types.Hotel
	objectId, _ := primitive.ObjectIDFromHex(id)
	if err := s.coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx *fasthttp.RequestCtx) ([]*types.Hotel, error) {
	var hotels []*types.Hotel
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}
