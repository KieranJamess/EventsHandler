package main

import (
	"github.com/KieranJamess/EventsHandler/config"
	"github.com/KieranJamess/EventsHandler/database"
	"github.com/KieranJamess/EventsHandler/shutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

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
	app := buildServer(env)

	// start the server
	go func() {
		app.Listen(":" + env.PORT)
	}()

	// return a function to close the server and database
	return func() {
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) *fiber.App {
	// create the fiber app
	app := fiber.New()

	// add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	//Connect to DB
	client, _ := database.ConnectToMongoDB(env.MONGO_URL)

	// Add event routes
	eventRoutes(app, client)

	return app
}
