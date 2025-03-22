package users

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"unit-test-mongo/internal/domain"
)

func (r repository) Get(ctx context.Context, userID string) (*domain.User, error) {
	// Convert string ID to ObjectID
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format '%s': %w", userID, err)
	}

	filter := bson.D{{"_id", objID}}

	var user domain.User
	err = r.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
