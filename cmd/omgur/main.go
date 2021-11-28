package main

import (
	"log"
	"net/http"
	"os"
	"git.geraldwu.com/gerald/omgur/pkg/app"
)

func main() {
	app, _ := app.CreateApp("Client-ID 546c25a59c58ad7")

	http.HandleFunc("/", app.HTTPServer)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/template/static"))))

	getenv:=func (key, fallback string) string {
	    value := os.Getenv(key)
	    if len(value) == 0 {
	        return fallback
	    }
	    return value
	}

	port:=getenv("PORT","8080")
	log.Print("Starting webserver on "+port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
