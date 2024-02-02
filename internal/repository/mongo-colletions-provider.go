package repository

import (
	"MongoDBCounterService/pkg/connections"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type CollectionsRepository interface {
	ListDatabaseNames(ctx context.Context) ([]string, error)
	ListCollectionNames(databaseName string, ctx context.Context) ([]string, error)
	LengthOfCollection(databaseName, collectionName string, ctx context.Context) (int64, error)
}

type MongoDB struct{}

func NewCollectionsRepository() CollectionsRepository {
	return &MongoDB{}
}

func (m *MongoDB) ListDatabaseNames(ctx context.Context) ([]string, error) {
	db := connections.GetMongoClient()
	return db.ListDatabaseNames(ctx, bson.M{})
}

func (m *MongoDB) ListCollectionNames(databaseName string, ctx context.Context) ([]string, error) {
	db := connections.GetMongoDB(databaseName)
	return db.ListCollectionNames(ctx, bson.M{})
}

func (m *MongoDB) LengthOfCollection(databaseName, collectionName string, ctx context.Context) (int64, error) {
	db := connections.GetMongoDB(databaseName)
	collection := db.Collection(collectionName)
	return collection.CountDocuments(ctx, bson.M{})
}
