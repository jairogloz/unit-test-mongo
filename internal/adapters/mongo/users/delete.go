package users

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r repository) Delete(ctx context.Context, userID string) error {
	// Convert string ID to ObjectID
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format '%s': %w", userID, err)
	}

	filter := bson.D{{"_id", objID}}

	result, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
