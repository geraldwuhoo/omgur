package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Image struct {
	Title       string
	Description string
	Link        string
}

type Album struct {
	Images []Image
}

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

	// Determine if this is an imgur.com/a/ album
	album, err := regexp.MatchString("a/.{7}", uri)
	if err != nil {
		log.Fatal(err)
	}

	if directImage {
		// This is a direct image, so use the direct image handler
		DirectImageHandler(w, uri)
	} else if album {
		// This is an album, so use the album handler
		AlbumHandler(w, uri)
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

	// If we were unable to get this image for any reason, then respond as such
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

func AlbumHandler(w http.ResponseWriter, uri string) {
	// Build GET request to Imgur API
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.imgur.com/3/album/%v/images", uri[2:]), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Client-ID 546c25a59c58ad7")

	// Execute GET request to get Album details
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// If we were unable to get this album for any reason, then respond as such
	if resp.StatusCode != 200 {
		output := fmt.Sprintf("Error %v looking up %v\n", resp.StatusCode, uri)
		log.Print(output)
		http.Error(w, output, resp.StatusCode)
		return
	}

	// Get contents of the API request
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Imgur response: %v\n", string(contents))

	// Unpack the JSON response into an unstructured Go struct
	var result map[string]interface{}
	json.Unmarshal([]byte(contents), &result)

	// Struct to store the album details in for templating
	data := Album{
		Images: []Image{},
	}
	// Loop over the results, and add each album image to the data struct
	for _, value := range result["data"].([]interface{}) {
		image := value.(map[string]interface{})

		var title string
		if image["title"] == nil {
			title = "No title"
		} else {
			title = image["title"].(string)
		}

		var description string
		if image["description"] == nil {
			description = "No description"
		} else {
			description = image["description"].(string)
		}

		link := image["link"].(string)
		link = link[strings.LastIndex(link, "/"):]

		log.Printf("link: %v\n", link)
		current_image := Image{
			Link:        link,
			Title:       title,
			Description: description,
		}
		data.Images = append(data.Images, current_image)
	}

	// Apply the extracted album to the template
	t, err := template.ParseFiles("templates/album.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
