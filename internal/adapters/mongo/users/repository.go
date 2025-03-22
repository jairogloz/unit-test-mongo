package users

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
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
