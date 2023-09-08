package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `json:"name"`
	StartDateTime time.Time          `json:"startDateTime"`
	EndDateTime   time.Time          `json:"endDateTime"`
}

func handleEvents(client *mongo.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		switch c.Method() {
		case fiber.MethodGet:
			// Handle GET /events
			handleGetEvents(c, client)

		case fiber.MethodPost:
			// Handle POST /events
			handlePostEvents(c, client)

		default:
			// Handle unsupported methods with an error response
			return fiber.NewError(fiber.StatusMethodNotAllowed, "Method not allowed")
		}
		return nil
	}
}

func handleGetEvents(c *fiber.Ctx, client *mongo.Client) error {
	// Function for handling GET method on /events
	var event Event
	err := handleGet(c, client, "events", 5*time.Second, &event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "GET event failed",
		})
	}

	return nil
}

func handlePostEvents(c *fiber.Ctx, client *mongo.Client) error {
	// Function for handling POST method on /events
	// Post event in events DB, events collection

	var event Event
	err := handlePost(c, client, "events", 5*time.Second, &event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Post event failed",
		})
	}

	return nil
}

func handleEventById(client *mongo.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		switch c.Method() {
		case fiber.MethodGet:
			// Handle GET /events
			handleGetEventsById(c, client)

		case fiber.MethodPatch:
			// Handle PATCH /events
			handlePatchEventsById(c, client)

		default:
			// Handle unsupported methods with an error response
			return fiber.NewError(fiber.StatusMethodNotAllowed, "Method not allowed")
		}
		return nil
	}
}

func handleGetEventsById(c *fiber.Ctx, client *mongo.Client) error {
	// Function for handling GET method on /events by ID
	var event Event
	err := handleGetByID(c, client, "events", 5*time.Second, &event, c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "GET event by ID failed",
		})
	}

	return nil
}

func handlePatchEventsById(c *fiber.Ctx, client *mongo.Client) error {
	// Function for handling PATCH method on /events by ID
	var event Event
	err := handlePatchById(c, client, "events", 5*time.Second, &event, c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "PATCH event by ID failed",
		})
	}
	return nil
}
