package main

import (
	"context"
	"reflect"
	"time"

	"github.com/KieranJamess/EventsHandler/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func eventRoutes(app *fiber.App, client *mongo.Client) {
	// Handle healthchecks
	app.Get("/healthcheck", healthCheckHandler)

	// Handle Events
	app.All("/events", handleEvents(client))
	app.All("/events/:id", handleEventById(client))
}

func handlePost(c *fiber.Ctx, client *mongo.Client, collectionName string, timeout time.Duration, format interface{}) error {
	// Parse the JSON request body into the  struct
	if err := c.BodyParser(format); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Insert the event into the MongoDB collection
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	insert, err := database.CreateCollection("events", collectionName, client).InsertOne(ctx, format)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert event",
		})
	}

	// Return a JSON response with the inserted event and a 201 Created status
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event successfully created",
		"id":      insert.InsertedID,
	})
}

func handleGet(c *fiber.Ctx, client *mongo.Client, collectionName string, timeout time.Duration, format interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{}

	cursor, err := database.CreateCollection("events", collectionName, client).Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve events",
		})
	}
	defer cursor.Close(ctx)
	results := []interface{}{}
	for cursor.Next(ctx) {
		item := reflect.New(reflect.TypeOf(format)).Interface()
		if err := cursor.Decode(item); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode event",
			})
		}
		results = append(results, item)
	}
	return c.JSON(results)
}

func handleGetByID(c *fiber.Ctx, client *mongo.Client, collectionName string, timeout time.Duration, format interface{}, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No ID specified",
		})
	}

	// Create a filter for retrieving a specific document by _id
	idString, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid ID String format",
		})
	}
	filter := bson.M{"_id": idString}

	cursor, err := database.CreateCollection("events", collectionName, client).Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve event by ID",
		})
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		item := reflect.New(reflect.TypeOf(format)).Interface()

		if err := cursor.Decode(item); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode event",
			})
		}

		return c.JSON(item)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Event not found",
	})
}

func handlePatchById(c *fiber.Ctx, client *mongo.Client, collectionName string, timeout time.Duration, format interface{}, id string) error {
	idString, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID string format",
		})
	}

	filter := bson.M{"_id": idString}

	if err := c.BodyParser(format); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create an update document with the fields you want to update
	update := bson.M{
		"$set": format,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Perform the update operation
	_, err = client.Database("events").Collection(collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update item in database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event successfully Updated",
		"id":      idString,
	})
}
