package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

//Ref: https://chatgpt.com/share/e3f542be-63c0-4051-95c8-edc47484d524

// MongoDBService defines the methods for CRUD operations.
type MongoDBService interface {
	InsertOne(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, collection string, filter interface{}) *mongo.SingleResult
	UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error)
}

// mongoDBService is the concrete implementation of MongoDBService.
type mongoDBService struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDBService creates a new MongoDBService.
func NewMongoDBService(client *mongo.Client, db *mongo.Database) MongoDBService {
	return &mongoDBService{
		client: client,
		db:     db,
	}
}

func (s *mongoDBService) InsertOne(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	return s.db.Collection(collection).InsertOne(ctx, document)
}

func (s *mongoDBService) FindOne(ctx context.Context, collection string, filter interface{}) *mongo.SingleResult {
	return s.db.Collection(collection).FindOne(ctx, filter)
}

func (s *mongoDBService) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return s.db.Collection(collection).UpdateOne(ctx, filter, update)
}

func (s *mongoDBService) DeleteOne(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	return s.db.Collection(collection).DeleteOne(ctx, filter)
}
