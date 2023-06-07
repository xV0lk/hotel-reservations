package db

import (
	"github.com/valyala/fasthttp"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserById(ctx *fasthttp.RequestCtx, id string) (*types.User, error)
	GetUsers(ctx *fasthttp.RequestCtx) ([]*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoUserStore) GetUserById(ctx *fasthttp.RequestCtx, id string) (*types.User, error) {
	var user types.User
	objectId, _ := primitive.ObjectIDFromHex(id)
	if err := s.coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx *fasthttp.RequestCtx) ([]*types.User, error) {
	var users []*types.User
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
