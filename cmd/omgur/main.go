package main

import (
	"log"
	"net/http"

	"git.geraldwu.com/gerald/omgur/pkg/app"
)

type Image struct {
	Title       string
	Description string
	Link        string
}

type Album struct {
	Title       string
	Description string
	Images      []Image
}

func main() {
	app, _ := app.CreateApp("Client-ID 546c25a59c58ad7")

	http.HandleFunc("/", app.HTTPServer)

	log.Print("Starting webserver on 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}