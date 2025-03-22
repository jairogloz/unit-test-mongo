package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
	"unit-test-mongo/internal/ports"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"unit-test-mongo/internal/adapters/mongo/users"
	"unit-test-mongo/internal/domain"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set client options - replace with your connection string if needed
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ensure disconnection when the function returns
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
		fmt.Println("Connection to MongoDB closed")
	}()

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	fmt.Println("Connected to MongoDB!")

	// Create a new user repository
	userRepo, err := users.NewRepository(client)
	if err != nil {
		log.Fatal("Failed to create user repository:", err)
	}

	// Run the test flow
	testRepositoryFlow(ctx, userRepo)
}

func testRepositoryFlow(ctx context.Context, repo ports.UserRepository) {
	fmt.Println("\n--- Starting Repository Test Flow ---")

	// Create a new user and print the insertedID
	user1ID, err := repo.Save(ctx, &domain.UserCreate{
		Name: "John Doe",
		Age:  30,
	})
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}
	fmt.Printf("Created user with ID: %s\n", user1ID)

	// Query the user by ID
	user1, err := repo.Get(ctx, user1ID)
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}
	fmt.Printf("Retrieved user: %+v\n", user1)

	// Modify the user
	user1.Name = "John Updated"
	user1.Age = 31
	err = repo.Update(ctx, user1ID, user1)
	if err != nil {
		log.Fatal("Failed to update user:", err)
	}
	fmt.Println("Updated user")

	// Query the updated user
	updatedUser, err := repo.Get(ctx, user1ID)
	if err != nil {
		log.Fatal("Failed to get updated user:", err)
	}
	fmt.Printf("Retrieved updated user: %+v\n", updatedUser)

	// Insert a second user
	user2ID, err := repo.Save(ctx, &domain.UserCreate{
		Name: "Jane Doe",
		Age:  28,
	})
	if err != nil {
		log.Fatal("Failed to create second user:", err)
	}
	fmt.Printf("Created second user with ID: %s\n", user2ID)

	// List both users
	userList, err := repo.List(ctx, 10, 0)
	if err != nil {
		log.Fatal("Failed to list users:", err)
	}
	fmt.Println("List of users:")
	for i, u := range userList {
		fmt.Printf("%d: %+v\n", i+1, u)
	}

	// Delete both users
	err = repo.Delete(ctx, user1ID)
	if err != nil {
		log.Fatal("Failed to delete first user:", err)
	}
	fmt.Printf("Deleted user with ID: %s\n", user1ID)

	// List users again (should be empty)
	emptyUsers, err := repo.List(ctx, 10, 0)
	if err != nil {
		log.Fatal("Failed to list users after deletion:", err)
	}
	fmt.Printf("Users after deletion: %d found\n", len(emptyUsers))

	fmt.Println("--- Repository Test Flow Completed ---")
}
