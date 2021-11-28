package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (a *App) GalleryHandler(w http.ResponseWriter, r *http.Request, uri string) {
	a.GetGallery(w, r, uri)
}

func (a *App) GetGallery(w http.ResponseWriter, r *http.Request, uri string) {
	endpoint := fmt.Sprintf("https://api.imgur.com/3/gallery/%v", uri[8:])
	// Build GET request to Imgur API
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", a.Authorization)

	// Execute GET request to get Gallery details
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// If we were unable to get this gallery for any reason, then respond as such
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

	// Unpack the JSON response into an unstructured Go struct
	var result map[string]interface{}
	json.Unmarshal([]byte(contents), &result)

	// Get parent data in JSON
	data := result["data"].(map[string]interface{})
	// Get if this gallery is an album or not
	var isAlbum bool
	if data["is_album"] == nil {
		log.Fatalf("This field can never be nil, fatal error looking up %v!", uri)
	} else {
		isAlbum = data["is_album"].(bool)
	}

	if isAlbum {
		// If this is an album, then we are good to just call get album directly
		log.Printf("%v is an album, calling GetAlbum", uri)
		a.GetAlbum(w, uri, endpoint)
	} else {
		// Otherwise, this is an actual image object, so call the imagehandler
		log.Printf("%v is an image, calling ImageHandler", uri)
		a.ImageHandler(w, r, uri[8:])
	}
}
