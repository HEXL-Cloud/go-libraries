package mongobase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/drivertest"
)

// TestEntity represents a test entity for testing purposes
type TestEntity struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func setupMockClient(t *testing.T, responses ...bson.D) (*mongo.Client, *drivertest.MockDeployment) {
	// Create a mock deployment with responses
	deployment := drivertest.NewMockDeployment(responses...)

	// Create a client with the mock deployment
	opts := options.Client()
	opts.Deployment = deployment
	client, err := mongo.Connect(opts)
	require.NoError(t, err)

	return client, deployment
}

func TestNew(t *testing.T) {
	client, _ := setupMockClient(t, bson.D{{Key: "ok", Value: 1}})
	defer client.Disconnect(context.Background())

	repo := New[TestEntity](client, "testdb", "testcol")

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.collection)
}

func TestMongoBaseRepository_InsertOne(t *testing.T) {
	tests := []struct {
		name          string
		document      TestEntity
		responses     []bson.D
		expectedError bool
	}{
		{
			name: "successful insertion",
			document: TestEntity{
				ID:   "test-id-1",
				Name: "John Doe",
				Age:  30,
			},
			responses: []bson.D{
				{
					{Key: "ok", Value: 1},
					{Key: "acknowledged", Value: true},
					{Key: "insertedId", Value: "test-id-1"},
				},
			},
			expectedError: false,
		},
		{
			name: "insertion error",
			document: TestEntity{
				ID:   "test-id-2",
				Name: "Jane Doe",
				Age:  25,
			},
			responses: []bson.D{
				{
					{Key: "ok", Value: 0},
					{Key: "errmsg", Value: "insertion failed"},
					{Key: "code", Value: 11000}, // Duplicate key error
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, deployment := setupMockClient(t)
			defer client.Disconnect(context.Background())

			deployment.AddResponses(tt.responses...)

			repo := New[TestEntity](client, "testdb", "testcol")
			err := repo.InsertOne(context.Background(), tt.document)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMongoBaseRepository_FindOneById(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		responses      []bson.D
		expectedResult TestEntity
		expectedError  bool
	}{
		{
			name: "successful find",
			id:   "test-id-1",
			responses: []bson.D{
				{
					{Key: "ok", Value: 1},
					{Key: "cursor", Value: bson.D{
						{Key: "id", Value: int64(0)},
						{Key: "ns", Value: "testdb.testcol"},
						{Key: "firstBatch", Value: bson.A{
							bson.D{
								{Key: "_id", Value: "test-id-1"},
								{Key: "name", Value: "John Doe"},
								{Key: "age", Value: 30},
							},
						}},
					}},
				},
			},
			expectedResult: TestEntity{
				ID:   "test-id-1",
				Name: "John Doe",
				Age:  30,
			},
			expectedError: false,
		},
		{
			name: "document not found",
			id:   "non-existent-id",
			responses: []bson.D{
				{
					{Key: "ok", Value: 1},
					{Key: "cursor", Value: bson.D{
						{Key: "id", Value: int64(0)},
						{Key: "ns", Value: "testdb.testcol"},
						{Key: "firstBatch", Value: bson.A{}},
					}},
				},
			},
			expectedResult: TestEntity{},
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, deployment := setupMockClient(t)
			defer client.Disconnect(context.Background())

			deployment.AddResponses(tt.responses...)

			repo := New[TestEntity](client, "testdb", "testcol")
			result, err := repo.FindOneById(context.Background(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestMongoBaseRepository_FindAll(t *testing.T) {
	tests := []struct {
		name           string
		filter         bson.M
		responses      []bson.D
		expectedResult []TestEntity
		expectedError  bool
	}{
		{
			name:   "successful find all",
			filter: bson.M{"age": bson.M{"$gte": 25}},
			responses: []bson.D{
				{
					{Key: "ok", Value: 1},
					{Key: "cursor", Value: bson.D{
						{Key: "id", Value: int64(0)},
						{Key: "ns", Value: "testdb.testcol"},
						{Key: "firstBatch", Value: bson.A{
							bson.D{
								{Key: "_id", Value: "test-id-1"},
								{Key: "name", Value: "John Doe"},
								{Key: "age", Value: 30},
							},
							bson.D{
								{Key: "_id", Value: "test-id-2"},
								{Key: "name", Value: "Jane Doe"},
								{Key: "age", Value: 25},
							},
						}},
					}},
				},
			},
			expectedResult: []TestEntity{
				{ID: "test-id-1", Name: "John Doe", Age: 30},
				{ID: "test-id-2", Name: "Jane Doe", Age: 25},
			},
			expectedError: false,
		},
		{
			name:   "empty result",
			filter: bson.M{"age": bson.M{"$gte": 100}},
			responses: []bson.D{
				{
					{Key: "ok", Value: 1},
					{Key: "cursor", Value: bson.D{
						{Key: "id", Value: int64(0)},
						{Key: "ns", Value: "testdb.testcol"},
						{Key: "firstBatch", Value: bson.A{}},
					}},
				},
			},
			expectedResult: nil, // Should be empty slice, not nil
			expectedError:  false,
		},
		{
			name:   "find error",
			filter: bson.M{"invalid": "query"},
			responses: []bson.D{
				{
					{Key: "ok", Value: 0},
					{Key: "errmsg", Value: "query failed"},
					{Key: "code", Value: 2}, // BadValue error
				},
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, deployment := setupMockClient(t)
			defer client.Disconnect(context.Background())

			deployment.AddResponses(tt.responses...)

			repo := New[TestEntity](client, "testdb", "testcol")
			result, err := repo.FindAll(context.Background(), tt.filter)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestMongoBaseRepository_UpdateOneById(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		update        bson.M
		responses     []bson.D
		expectedError bool
	}{
		{
			name:   "successful update",
			id:     "test-id-1",
			update: bson.M{"name": "Updated Name", "age": 35},
			responses: []bson.D{
				{
					{Key: "ok", Value: 1},
					{Key: "acknowledged", Value: true},
					{Key: "matchedCount", Value: 1},
					{Key: "modifiedCount", Value: 1},
				},
			},
			expectedError: false,
		},
		{
			name:   "update error",
			id:     "test-id-2",
			update: bson.M{"invalid": "update"},
			responses: []bson.D{
				{
					{Key: "ok", Value: 0},
					{Key: "errmsg", Value: "update failed"},
					{Key: "code", Value: 2}, // BadValue error
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, deployment := setupMockClient(t)
			defer client.Disconnect(context.Background())

			deployment.AddResponses(tt.responses...)

			repo := New[TestEntity](client, "testdb", "testcol")
			err := repo.UpdateOneById(context.Background(), tt.id, tt.update)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMongoBaseRepository_DeleteOneById(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		responses     []bson.D
		expectedError bool
	}{
		{
			name: "successful deletion",
			id:   "test-id-1",
			responses: []bson.D{
				{
					{Key: "ok", Value: 1},
					{Key: "acknowledged", Value: true},
					{Key: "deletedCount", Value: 1},
				},
			},
			expectedError: false,
		},
		{
			name: "deletion error",
			id:   "test-id-2",
			responses: []bson.D{
				{
					{Key: "ok", Value: 0},
					{Key: "errmsg", Value: "delete failed"},
					{Key: "code", Value: 2}, // BadValue error
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, deployment := setupMockClient(t)
			defer client.Disconnect(context.Background())

			deployment.AddResponses(tt.responses...)

			repo := New[TestEntity](client, "testdb", "testcol")
			err := repo.DeleteOneById(context.Background(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
