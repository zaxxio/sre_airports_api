package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type Airport struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	IATA     string `json:"iata"`
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

const base64Creds = "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAic3JlLWFpcnBvcnRzLWFwaSIsCiAgInByaXZhdGVfa2V5X2lkIjogImI0Y2U0OTVmMDdmY2RmN2Y3ZGM3M2ZkZGY1MTc1YTI0NGYwOTEwOTIiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQ28zWWpwTnNRazFCZTVcbkxTMklueFkyWGwrMUpWVWJPS2tNL2w4RlpEeHg4S2t0U1M5N2Y1Y3BFNFk4YXlHZVJUeEZOS1pEZjlyNFdub1lcbnBtMDBsb2dDT3lLYnY2Nmo4cGx1ZzhRT2IyKzJ2NkVUK0h0U0x2YVNsM09lcUxxOTQ5OE1Jd2ZDTUc0dUYyWkdcbkdFRjQ5MUtiWjRxSTNrZGVaZG9XdVFaQVVTZHJEcjl1VlJZYjZ1Tjd4WkhscHNWYnZtVktiTDNwT1FRUEVEOWdcbm0yYUtvUURQWnZmOG1nQVVJUzRHUkwyZ0xYTWRFdFFRVzJ1dlRPTGFOYTRONlFyVVFCWWJoSVZVVnp0czlITk5cbnc4Q1hNc0xHbTFyYkJES05UU3dxcFpTekxuWlZkZHJKMXl6Zzc5VzJVUlAwdXRPODl3QVd0M1Bvd2l2RHlRYkpcbnZ1NStHSDBWQWdNQkFBRUNnZ0VBQUpPT20yRXVrNTAwaExtWURYNmhoUlF3WXo0ejRSRE92dW5va0IrbDl5dVpcbk9kVS9CTDdNUHpDL0s2Y0dYS2JXMlNCcyt0d3o5cU9ubnp5RWw0SkMxR01lcjlKcExPVVdNZng1d1NUY0NFSmJcbnlySW92QUNVYjN0MTNFVjYrQm1QNjZTVnhINmUwMElNL1RnUHhzQUc5Nmc2dUZONnRSV3VIY3E5U1orRHNRZkFcbnZIOWJ5ZUZVOGxMbkJVaDR1Y2dBdS9mZU5nNWhBbHU1WVpEbXRlQVdBTW9hUyt4QjVBanFyMzd3d0s3RVY5YTdcbklrU3Mxc0gzOVZDRkJqL2wrLzhBbUZiSUVqMkw0VUtFK25CdUxtZ2h5dndZNlZkUGJhczN1UGhmc0VOcTZZQXFcbkFhalhPSllkU1JCTTJQcDFJWitBdG1scG12WnNUeDBMUzlNT3ZJRTJad0tCZ1FEWDdhMGFBdDF5dzFoZkhsTjNcbi8yVEoveU1yM0pVWThPcFo0VHJNYkYvbXA5REgvczFEMlV1YmY3eHBQanU2TUxiNEtDQWI3cXdFSmRNbFJIYmxcbi8wSTRpbFFJZXVjcG5JV080MmdCVkxpOUpMeEFZaEN3SGRZbEtIYzBhS2IxU2lSekt6ZWdLYUN4ekE5cnFLVE5cblc2bGdqTHJtaUlmYVU4cW5vcFFYS2UzSEN3S0JnUURJTS96cUZpc01veFczVU1hSVBJRFhhRm93RDlzZUhXeFlcbmhSNlpBQ2sxaUovUStrNEpKRWY4OEhrckN5MVF4dmVHL29sM0ZQWWIxWmoyZEdSQ0JTK0t4ODIzRXQydHlaRkxcbmRUaGNyMmhJeWUrempFUTVWVFNwbytNOHFIb0hRcDZIU3N6R2Vnc3lJVzY3eHNJK2w0SWNKWHI4M2pGUFZmU0VcbjF1SlU2VmJnWHdLQmdEOHhnL09VMnhKM01TbkZTbEJZSWpzcnZETmQveFNwalN4NHlpaUJueDkyQlpoQ2JmaHBcblk4TkNndldhRFFqVXNQZTNabzVHTDNtWFNGQWoxVmhDZURMcjZPUUNkQnl0ZmpqdlBNVUc4bm9JZ2orbGM1VFhcblpwREJZd0duanhWQ2VhQnJDWUNLTGtsYW16aTZ4bUNEYnZLZXZTUXkyTytBamxLNU5mWUJnMkU1QW9HQUV1VGRcbkpKWnMvNmRRZ0ZsdU15TktvWW1tb1V5TnlGek1nZG9tVmhndXkyK1diWm1CemRrUHRpNVhzUmsvOEpTbWZhWCtcbkFUQUlQZjQ5amx6VHJXdGgzajRYQ3dVTHlML3lKMlhycU11aEV1V0Q2clQ5SjFBRVJWSkRPdEZIbXZITmxrVVhcbjZFOVNTU3Zna0hZa2xOV2xvTlJrdEFLZ01yV1Erd3h3bGNUanZ3OENnWUVBeTlvTk9wVmh3VDI4anJQUmp0ZVdcbmlGWk9jV2FVQmhFUEhMSlN4aDhyRHlGV3ZVbmdHVWd5bHB0OUh0aHA3dE1LUDZtK1dmenFwR0U2TFZSdDhJOEhcbmFPQStra2RVWlRCMlNIa0FIemFHWmFWOXZZYmFVV3N5VGprck80dmR5cllDa2FWdG1zOG5mRjFFYy8vZEhXYWlcbkVVTkxoY3dNZ29VUnhXekVTUE1JOFowPVxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogInRlcnJhZm9ybS1zZXJ2aWNlLWFjY291bnRAc3JlLWFpcnBvcnRzLWFwaS5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIsCiAgImNsaWVudF9pZCI6ICIxMDIwMzA4Mzc4ODc5NjQ1MDIyNzIiLAogICJhdXRoX3VyaSI6ICJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9vYXV0aDIvYXV0aCIsCiAgInRva2VuX3VyaSI6ICJodHRwczovL29hdXRoMi5nb29nbGVhcGlzLmNvbS90b2tlbiIsCiAgImF1dGhfcHJvdmlkZXJfeDUwOV9jZXJ0X3VybCI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHMiLAogICJjbGllbnRfeDUwOV9jZXJ0X3VybCI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9yb2JvdC92MS9tZXRhZGF0YS94NTA5L3RlcnJhZm9ybS1zZXJ2aWNlLWFjY291bnQlNDBzcmUtYWlycG9ydHMtYXBpLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAidW5pdmVyc2VfZG9tYWluIjogImdvb2dsZWFwaXMuY29tIgp9Cg=="

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

	decodedCreds, err := base64.StdEncoding.DecodeString(base64Creds)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding Base64: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the decoded credentials to a temporary file
	tempFile, err := ioutil.TempFile("", "gcs-credentials-*.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp file: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name()) // Clean up the file after we're done

	if _, err := tempFile.Write(decodedCreds); err != nil {
		http.Error(w, fmt.Sprintf("Error writing to temp file: %v", err), http.StatusInternalServerError)
		return
	}
	if err := tempFile.Close(); err != nil {
		http.Error(w, fmt.Sprintf("Error closing temp file: %v", err), http.StatusInternalServerError)
		return
	}

	// Initialize GCS client
	creds := option.WithCredentialsFile(tempFile.Name())
	ctx := context.Background()
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		http.Error(w, "Failed to create GCS client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Define GCS bucket name and object name (file path)
	bucketName := "airportima-bucket"
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
	http.HandleFunc("/update_airport_image", UpdateAirportImage)

	// Start the server
	http.ListenAndServe(":8080", nil)
}
