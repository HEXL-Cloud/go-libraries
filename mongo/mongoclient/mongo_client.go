package mongoclient

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
}

// Creates a new MongoClient instance and checks the connection by pinging the server.
//
// Parameters:
//   - ctx: The context for the operation.
//   - connectionStr: The MongoDB connection string.
//
// Returns:
//   - A pointer to the MongoClient instance or an error if the connection fails.
func Connect(ctx context.Context, connectionStr string) (*MongoClient, error) {
	// Create a new MongoDB client
	client, err := mongo.Connect(options.Client().ApplyURI(connectionStr))
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoClient{
		Client: client,
	}, nil
}

// Closes the MongoDB client connection.
//
// Parameters:
//   - ctx: The context for the operation.
//
// Returns:
//   - An error if the disconnection fails.
func (mc *MongoClient) Disconnect(ctx context.Context) error {
	if err := mc.Client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
