package users

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"unit-test-mongo/internal/domain"
	"unit-test-mongo/internal/ports"
)

// repository implements ports.UserRepository
type repository struct {
	coll *mongo.Collection
}

func NewRepository(client *mongo.Client) (ports.UserRepository, error) {
	r := &repository{
		coll: client.Database("test").Collection("users"),
	}
	return r, nil
}

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

func (r repository) Update(ctx context.Context, userID string, user *domain.User) error {
	// Convert string ID to ObjectID
	objID, err := bson.ObjectIDFromHex(userID)
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

func (r repository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{"_id", 1}})

	cursor, err := r.coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
