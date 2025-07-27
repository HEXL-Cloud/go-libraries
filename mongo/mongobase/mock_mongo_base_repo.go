package mongobase

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MockMongoBaseRepository[T any] struct {
	mock.Mock
}

func (m *MockMongoBaseRepository[T]) InsertOne(ctx context.Context, document T) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *MockMongoBaseRepository[T]) FindOneById(ctx context.Context, id string) (T, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return *new(T), args.Error(1)
	}
	return args.Get(0).(T), args.Error(1)
}

func (m *MockMongoBaseRepository[T]) FindAll(ctx context.Context, filter bson.M) ([]T, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockMongoBaseRepository[T]) UpdateOneById(ctx context.Context, id string, update bson.M) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockMongoBaseRepository[T]) DeleteOneById(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
