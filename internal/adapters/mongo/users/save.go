package users

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"unit-test-mongo/internal/domain"
)

func (r repository) Save(ctx context.Context, user *domain.UserCreate) (string, error) {
	result, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	// Extract the inserted ID and convert it to string
	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		// Convert to string if it's not already a string type
		return fmt.Sprintf("%v", result.InsertedID), nil
	}

	return insertedID.Hex(), nil
}
