package users_test

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"unit-test-mongo/internal/adapters/mongo/users"

	"github.com/stretchr/testify/assert"
	"unit-test-mongo/internal/domain"
)

func TestRepository_Save(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	t.Run("success", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			// Setup mock response
			mt.AddMockResponses(mtest.CreateSuccessResponse(
				bson.E{Key: "n", Value: 1},
			))

			repo, err := users.NewRepository(mt.Client)
			if err != nil {
				mt.Errorf("error creating repository: %v", err)
			}

			// Create test user
			user := &domain.UserCreate{
				Name: "Test User",
				Age:  30,
			}

			// Call the function
			id, err := repo.Save(context.Background(), user)

			// Assert results
			assert.NoError(t, err)
			assert.NotEmpty(t, id)
		})
	})

	t.Run("database error", func(t *testing.T) {
		mt.Run("database error", func(mt *mtest.T) {
			// Setup mock error response
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "duplicate key error",
			}))

			repo, err := users.NewRepository(mt.Client)
			if err != nil {
				mt.Errorf("error creating repository: %v", err)
			}

			// Create test user
			user := &domain.UserCreate{
				Name: "Test User",
				Age:  30,
			}

			// Call the function
			id, err := repo.Save(context.Background(), user)

			// Assert results
			assert.Error(t, err)
			assert.Equal(t, "", id)
		})
	})

}
