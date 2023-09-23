package db

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userColl = "users"

type Dropper interface {
	Drop(context.Context) error
}
type UserStore interface {
	IndexEmail(ctx context.Context) error
	GetUserById(ctx *fasthttp.RequestCtx, id string) (*types.User, error)
	GetUsers(ctx *fasthttp.RequestCtx) ([]*types.User, error)
	GetUser(ctx *fasthttp.RequestCtx, filter bson.M) (*types.User, error)
	InsertUser(ctx context.Context, user *types.User) error
	DeleteUser(ctx *fasthttp.RequestCtx, id string) error
	UpdateUser(ctx *fasthttp.RequestCtx, id string, updateUser *types.UpdateUserParams) (*types.User, error)

	Dropper
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, dbname string) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoUserStore) IndexEmail(ctx context.Context) error {
	_, err := s.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("Dropping test database")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) GetUserById(ctx *fasthttp.RequestCtx, id string) (*types.User, error) {
	var user types.User
	objectId, _ := primitive.ObjectIDFromHex(id)
	if err := s.coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUser(ctx *fasthttp.RequestCtx, filter bson.M) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, filter).Decode(&user); err != nil {
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

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) error {
	result, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *MongoUserStore) DeleteUser(ctx *fasthttp.RequestCtx, id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx *fasthttp.RequestCtx, id string, updateUser *types.UpdateUserParams) (*types.User, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	// update := bson.M{"$set": updateUser}
	update := bson.M{"$set": updateUser.ToBson()}
	result, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	if updated, err := s.GetUserById(ctx, id); err != nil {
		return nil, err
	} else {
		return updated, nil
	}
}
