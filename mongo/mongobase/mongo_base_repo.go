package mongobase

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// This interface defines the list of methods currently supported by the MongoBaseRepository.
// It should not be used directly in the application, but rather as a reference for implementing specific repositories.
//
// Any method not listed here can be implemented as an extension.
type IMongoBaseRepository interface {
	InsertOne(ctx context.Context, document Entity) error
	FindOneById(ctx context.Context, id string) (Entity, error)
	FindAll(ctx context.Context, filter bson.M) ([]Entity, error)
	UpdateOneById(ctx context.Context, id string, update bson.M) error
	DeleteOneById(ctx context.Context, id string) error
}

type MongoBaseRepository[T any] struct {
	collection *mongo.Collection
}

// Creates a new instance of the MongoBaseRepository
// Parameters:
//   - client: The MongoDB client to use
//   - databaseName: The name of the database
//   - collectionName: The name of the collection to operate on
//
// Returns:
//   - A pointer to the MongoBaseRepository instance
//
// Note:
//   - This contains the implementation for base operations. It can be extended for specific entity repositories if required.
//
// Example usage:
//
//	userRepository := mongobase.New[entity.User](client, "foo-db", "users")
//
// Extending the usage:
//
//	type UserRepository struct {
//		MongoBaseRepository[entity.User]
//	}
//
//	func (UserRepository) CreateUserWithCustomLogic(ctx context.Context, user entity.User) error { ... }
//	// ... other methods specific to UserRepository
func New[T any](client *mongo.Client, databaseName, collectionName string) *MongoBaseRepository[T] {
	collection := client.Database(databaseName).Collection(collectionName)
	return &MongoBaseRepository[T]{
		collection: collection,
	}
}

// Inserts a single document into the collection
//
// Parameters:
//   - ctx: The context for the operation
//   - document: The document to insert
//
// Returns:
//   - nil if the insertion is successful
//   - An error if the insertion fails
func (repo *MongoBaseRepository[T]) InsertOne(ctx context.Context, document T) error {
	_, err := repo.collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

// Finds a single document by _id
//
// Parameters:
//   - ctx: The context for the operation
//   - id: The ID of the document to find
//
// Returns:
//   - The found document if successful
//   - An error if the document is not found or if there is an error during the operation
func (repo *MongoBaseRepository[T]) FindOneById(ctx context.Context, id string) (T, error) {
	filter := bson.M{"_id": id}

	var result T
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Finds all documents matching the provided filter
//
// Parameters:
//   - ctx: The context for the operation
//   - filter: The filter to apply to the query
//
// Returns:
//   - A slice of found documents of type T
//   - An error if the operation fails
func (repo *MongoBaseRepository[T]) FindAll(ctx context.Context, filter bson.M) ([]T, error) {
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var item T
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Updates a single document by _id
//
// Parameters:
//   - ctx: The context for the operation
//   - id: The ID of the document to update
//   - update: The update to apply to the document
//
// Returns:
//   - nil if the update is successful
func (repo *MongoBaseRepository[T]) UpdateOneById(ctx context.Context, id string, update bson.M) error {
	filter := bson.M{"_id": id}
	_, err := repo.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return err
	}

	return nil
}

// Deletes a single document by _id
//
// Parameters:
//   - ctx: The context for the operation
//   - id: The ID of the document to delete
//
// Returns:
//   - nil if the deletion is successful
func (repo *MongoBaseRepository[T]) DeleteOneById(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
