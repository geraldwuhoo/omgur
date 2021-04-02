package app

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

type App struct {
	Authorization string
	Rdb           *redis.Client
}

func CreateApp(authorization string) (*App, error) {
	// Create new app struct
	app := new(App)

	// Set the authorization client ID
	app.Authorization = authorization

	// Set the Redis host connection details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort, exists := os.LookupEnv("REDIS_PORT")
	if !exists {
		redisPort = "6379"
	}

	// Create the Redis connection (if applicable)
	if redisHost != "" {
		app.Rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", redisHost, redisPort),
			Password: "",
			DB:       0,
		})
	}

	return app, nil
}
