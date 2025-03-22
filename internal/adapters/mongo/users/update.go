package users

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"unit-test-mongo/internal/domain"
)

func (r repository) Update(ctx context.Context, userID string, user *domain.User) error {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format '%s': %w", userID, err)
	}

	filter := bson.D{{"_id", objID}}

	// Create update document - exclude _id from the update
	update := bson.D{{"$set", bson.D{
		{"name", user.Name},
		{"age", user.Age},
	}}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
