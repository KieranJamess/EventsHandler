package main

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func eventRoutes(app *fiber.App, client *mongo.Client) {
	// Handle healthchecks
	app.Get("/healthcheck", healthCheckHandler)

	// Handle Events
	app.All("/events", eventHandler(client))
}
