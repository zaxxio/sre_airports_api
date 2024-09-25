package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type Airport struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	IATA    string `json:"iata"`
	ImageURL string `json:"image_url"`
}

type AirportV2 struct {
	Airport
	RunwayLength int `json:"runway_length"`
}

// Mock data for airports in Bangladesh
var airports = []Airport{
	{"Hazrat Shahjalal International Airport", "Dhaka", "DAC", "https://storage.googleapis.com/bd-airport-data/dac.jpg"},
	{"Shah Amanat International Airport", "Chittagong", "CGP", "https://storage.googleapis.com/bd-airport-data/cgp.jpg"},
	{"Osmani International Airport", "Sylhet", "ZYL", "https://storage.googleapis.com/bd-airport-data/zyl.jpg"},
}

// Mock data for airports in Bangladesh (with runway length for V2)
var airportsV2 = []AirportV2{
	{Airport{"Hazrat Shahjalal International Airport", "Dhaka", "DAC", "https://storage.googleapis.com/bd-airport-data/dac.jpg"}, 3200},
	{Airport{"Shah Amanat International Airport", "Chittagong", "CGP", "https://storage.googleapis.com/bd-airport-data/cgp.jpg"}, 2900},
	{Airport{"Osmani International Airport", "Sylhet", "ZYL", "https://storage.googleapis.com/bd-airport-data/zyl.jpg"}, 2500},
}

// HomePage handler
func HomePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: OK"))
}

// Airports handler for the first endpoint
func Airports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airports)
}

// AirportsV2 handler for the second version endpoint
func AirportsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airportsV2)
}

// ##############################
// ## TODO: Edit this function ##
// ##############################

// UpdateAirportImage handler for updating airport images
// UpdateAirportImage handler for updating airport images using airport name
func UpdateAirportImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the multipart form data (limit to 10 MB)
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve file from form
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the airport name from the form
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Airport name is required", http.StatusBadRequest)
		return
	}

	// Initialize GCS client
	creds := option.WithCredentialsFile("./creds/dummy.json")
	ctx := context.Background()
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		http.Error(w, "Failed to create GCS client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Define GCS bucket name and object name (file path)
	bucketName := "bd-airport-data"
	objectName := fmt.Sprintf("%s-%s", name, handler.Filename)

	// Upload image to GCS
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	wc := object.NewWriter(ctx)

	if _, err = io.Copy(wc, file); err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}
	if err := wc.Close(); err != nil {
		http.Error(w, "Error closing writer", http.StatusInternalServerError)
		return
	}

	// Construct the GCS URL
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	// Update the in-memory airport data by matching the name
	updated := false
	for i := range airports {
		if airports[i].Name == name {
			airports[i].ImageURL = imageURL
			updated = true
			break
		}
	}

	if !updated {
		http.Error(w, "Airport not found", http.StatusNotFound)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Image uploaded successfully",
		"image_url": imageURL,
	})
}

func main() {
	// Setup routes
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/airports", Airports)
	http.HandleFunc("/airports_v2", AirportsV2)

	// TODO: complete the UpdateAirportImage handler function
	http.HandleFunc("/update_airport_image", UpdateAirportImage)

	// Start the server
	http.ListenAndServe(":8080", nil)
}
