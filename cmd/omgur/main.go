package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"git.geraldwu.com/gerald/omgur/pkg/app"
)

func main() {
	app, _ := app.CreateApp("Client-ID 546c25a59c58ad7")
	fsStatic, err := fs.Sub(app.Content, "web/template/static")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", app.HTTPServer)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fsStatic))))

	port, ok := os.LookupEnv("OMGUR_LISTEN_PORT")
	if !ok {
		port = "8080"
	}
	
	address, ok := os.LookupEnv("OMGUR_LISTEN_ADDRESS")
        if !ok {
                address = ""
        }
	
	

	log.Printf("Starting webserver on %v:%v", address, port)
	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", address, port), nil); err != nil {
		log.Fatal(err)
	}
}
