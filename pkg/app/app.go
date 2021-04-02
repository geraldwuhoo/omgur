package app

type App struct {
	Authorization string
}

func CreateApp(authorization string) (*App, error) {
	app := new(App)
	app.Authorization = authorization

	return app, nil
}
