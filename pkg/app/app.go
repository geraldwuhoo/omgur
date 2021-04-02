package app

import (
	"github.com/go-redis/redis/v8"
)

type App struct {
	Authorization string
	Rdb           *redis.Client
}

func CreateApp(authorization string) (*App, error) {
	app := new(App)
	app.Authorization = authorization

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	app.Rdb = rdb

	return app, nil
}
