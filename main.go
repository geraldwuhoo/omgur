package main

import (
    "net/http"
    "log"
    "fmt"
    "io"
    "regexp"
)

func main() {
    http.HandleFunc("/", HTTPServer)

    log.Print("Starting webserver on 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

func HTTPServer(w http.ResponseWriter, r *http.Request) {
    log.Printf("Serving page for %s", r.URL.Path)
    uri := r.URL.Path[1:]

    // Security
    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }

    // Determine if this is a direct i.imgur.com image
    directImage, err := regexp.MatchString(".*\\.(jpg|jpeg|png|gif|gifv|apng|tiff|mp4|mpeg|avi|webm)", uri)
    if err != nil {
        log.Print(err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    if directImage {
        // This is a direct image, so use the direct image handler
        DirectImageHandler(w, uri)
    } else {
        // Future proxying features not yet implemented
        http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
        return
    }
}

func DirectImageHandler(w http.ResponseWriter, uri string) {
    // Get the image directly from i.imgur.com
    resp, err := http.Get(fmt.Sprintf("https://i.imgur.com/%v", uri))
    if err != nil {
        log.Print(err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // If we were unable to get this iamge for any reason, then respond as such
    if resp.StatusCode != 200 {
        output := fmt.Sprintf("Error %v looking up %v\n", resp.StatusCode, uri)
        log.Print(output)
        http.Error(w, output, resp.StatusCode)
        return
    }

    // Successfully got image, so return the proper response as a direct image
    w.Header().Set("Content-Length", fmt.Sprint(resp.ContentLength))
    w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

    if _, err = io.Copy(w, resp.Body); err != nil {
        log.Print(err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
