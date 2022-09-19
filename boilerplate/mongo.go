package boilerplate

import (
	"context"
	"errors"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	MongoURI        string
	mongoClient     *mongo.Client
	MongoClientOnce sync.Once
}

type MongoCollection struct {
	MongoClient *MongoClient
	Database    string
	Collection  string
}

// Validate mongo configs
func (m *MongoClient) Validate() error {
	if m.MongoURI == "" {
		return errors.New("Missing mongo uri")
	}
	return nil
}

// GetSetupMongoClient get mongo client
func (m *MongoClient) GetSetupMongoClient() *mongo.Client {
	m.MongoClientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(m.MongoURI)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatalf("Error while connecting to mongo: %v", err)
		}
		m.mongoClient = client

	})
	return m.mongoClient
}

// GetCollection get mongo collection
func (m *MongoCollection) GetCollectionRef() *mongo.Collection {
	return m.MongoClient.GetSetupMongoClient().Database(m.Database).Collection(m.Collection)
}
