package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
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

func (a *App) AlbumHandler(w http.ResponseWriter, uri string) {
	// Build GET request to Imgur API
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.imgur.com/3/album/%v", uri[2:]), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", a.Authorization)

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

	album := a.ParseAlbum(contents)

	// Apply the extracted album to the template
	t, err := template.ParseFiles("web/template/album.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, *album)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) ParseAlbum(contents []byte) *Album {

	// Unpack the JSON response into an unstructured Go struct
	var result map[string]interface{}
	json.Unmarshal([]byte(contents), &result)

	// Get parent data in JSON
	data := result["data"].(map[string]interface{})
	// Get post title (safely)
	var postTitle string
	if data["title"] == nil {
		postTitle = ""
	} else {
		postTitle = data["title"].(string)
	}
	// Get post description (safely)
	var postDesc string
	if data["description"] == nil {
		postDesc = ""
	} else {
		postDesc = data["description"].(string)
	}
	// Get images slice
	images := data["images"].([]interface{})

	// Struct to store the album details in for templating
	album := &Album{
		Title:       postTitle,
		Description: postDesc,
		Images:      []Image{},
	}
	// Loop over the results, and add each album image to the data struct
	for _, value := range images {
		// Assert type of the image JSON data
		image := value.(map[string]interface{})

		// Get the title (safely)
		var title string
		if image["title"] == nil {
			title = ""
		} else {
			title = image["title"].(string)
		}

		// Get the description (safely)
		var description string
		if image["description"] == nil {
			description = ""
		} else {
			description = image["description"].(string)
		}

		// Get the link and parse the uri
		link := image["link"].(string)
		link = link[strings.LastIndex(link, "/"):]
		log.Printf("link: %v\n", link)

		// Add this image to the overall data
		album.Images = append(album.Images, Image{
			Link:        link,
			Title:       title,
			Description: description,
		})
	}

	return album
}
