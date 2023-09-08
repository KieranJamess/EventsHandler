package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event struct {
	Name          string    `json:"name"`
	StartDateTime time.Time `json:"startDateTime"`
	EndDateTime   time.Time `json:"endDateTime"`
}

func eventHandler(client *mongo.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		switch c.Method() {
		case fiber.MethodGet:
			// Handle GET /events
			getEventHandler(c, client)

		case fiber.MethodPost:
			// Handle POST /events
			posttEventHandler(c, client)

		case fiber.MethodDelete:
			// Handle DELETE /events
			deleteEventHandler(c, client)

		default:
			// Handle unsupported methods with an error response
			return fiber.NewError(fiber.StatusMethodNotAllowed, "Method not allowed")
		}
		return nil
	}
}

func getEventHandler(c *fiber.Ctx, client *mongo.Client) error {
	// Function for handling GET method on /events
	// Get all events in BSON format from events mongo db table
	return nil
}

func posttEventHandler(c *fiber.Ctx, client *mongo.Client) error {
	// Function for handling POST method on /events
	// Post event in events DB

	ctx := context.TODO()

	databaseName := "events"
	collectionName := "events"
	collection := client.Database(databaseName).Collection(collectionName)

	var event Event
	if err := c.BodyParser(&event); err != nil {
		return err
	}

	// Insert the event into the collection
	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		return err
	}

	return c.SendString("Event successfully inserted into the database")
}

func deleteEventHandler(c *fiber.Ctx, client *mongo.Client) error {
	// Function for handling DELETE method on /events
	return nil
}
