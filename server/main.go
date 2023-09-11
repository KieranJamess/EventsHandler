package main

import (
	"github.com/KieranJamess/EventsHandler/config"
	"github.com/KieranJamess/EventsHandler/database"
	"github.com/KieranJamess/EventsHandler/shutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerComponents struct {
	App    *fiber.App
	Client *mongo.Client
}

func main() {
	// load config
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	cleanup, err := run(env)

	defer cleanup()

	if err != nil {
		panic(err)
	}

	shutdown.Shutdown()
}

func run(env config.EnvVars) (func(), error) {
	components := buildServer(env)
	app := components.App

	// start the server
	go func() {
		app.Listen(":" + env.PORT)
	}()

	// return a function to close the server and database
	return func() {
		database.CloseMongoDBConnection(components.Client)
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) ServerComponents {
	// Create the fiber app
	app := fiber.New()

	// Add middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(logger.New())

	// Connect to DB
	client, _ := database.ConnectToMongoDB(env.MONGO_URL)

	// Add event routes
	eventRoutes(app, client)

	return ServerComponents{
		App:    app,
		Client: client,
	}
}
