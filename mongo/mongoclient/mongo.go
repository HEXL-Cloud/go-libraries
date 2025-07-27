package mongoclient

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

// Creates a new MongoClient instance and checks the connection by pinging the server.
// Parameters:
//   - connectionStr: The MongoDB connection string.
//
// Returns:
//   - A pointer to the MongoClient instance or an error if the connection fails.
func Connect(connectionStr string) (*MongoClient, error) {
	// Create a new MongoDB client
	client, err := mongo.Connect(options.Client().ApplyURI(connectionStr))
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return &MongoClient{
		client: client,
	}, nil
}

// Closes the MongoDB client connection.
// Returns:
//   - An error if the disconnection fails.
func (mc *MongoClient) Disconnect() error {
	if err := mc.client.Disconnect(context.Background()); err != nil {
		return err
	}

	return nil
}
