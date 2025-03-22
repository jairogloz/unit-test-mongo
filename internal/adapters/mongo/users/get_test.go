package users_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"unit-test-mongo/internal/adapters/mongo/users"
	"unit-test-mongo/internal/domain"
)

func TestRepository_Get(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	t.Run("success", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			// Create valid ObjectID
			id := primitive.NewObjectID()

			// Setup mock response - a document matching the user
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{
				{Key: "_id", Value: id},
				{Key: "name", Value: "Test User"},
				{Key: "age", Value: 30},
			}))

			repo, err := users.NewRepository(mt.Client)
			if err != nil {
				mt.Errorf("error creating repository: %v", err)
			}

			// Call the function with the valid ID
			user, err := repo.Get(context.Background(), id.Hex())

			// Assert results
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, "Test User", user.Name)
			assert.Equal(t, 30, user.Age)
		})
	})

	t.Run("user not found", func(t *testing.T) {
		mt.Run("user not found", func(mt *mtest.T) {
			// Setup mock response for no document found
			mt.AddMockResponses(mtest.CreateCursorResponse(
				0,
				"test.users",
				mtest.FirstBatch, // Empty batch
			))

			repo, err := users.NewRepository(mt.Client)
			if err != nil {
				mt.Errorf("error creating repository: %v", err)
			}

			id := primitive.NewObjectID()
			user, err := repo.Get(context.Background(), id.Hex())

			if assert.Error(t, err) {
				assert.ErrorIs(t, err, domain.ErrNotFound)
			}

			assert.Nil(t, user)
		})
	})

	t.Run("database error", func(t *testing.T) {
		mt.Run("database error", func(mt *mtest.T) {
			// Setup mock error response
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    12345,
					Message: "database error",
				},
			))

			repo, err := users.NewRepository(mt.Client)
			if err != nil {
				mt.Errorf("error creating repository: %v", err)
			}

			// Call the function with a valid ObjectID format
			id := primitive.NewObjectID()
			user, err := repo.Get(context.Background(), id.Hex())

			// Assert results
			assert.Error(t, err)
			assert.Nil(t, user)
		})
	})

	t.Run("invalid id format", func(t *testing.T) {
		mt.Run("invalid id format", func(mt *mtest.T) {
			repo, err := users.NewRepository(mt.Client)
			if err != nil {
				mt.Errorf("error creating repository: %v", err)
			}

			// Call the function with an invalid ObjectID format
			user, err := repo.Get(context.Background(), "invalid-id")

			// Assert results
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Contains(t, err.Error(), "invalid user ID format")
		})
	})
}
