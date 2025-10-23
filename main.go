package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type GeocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
	} `json:"results"`
}

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	CurrentWeather struct {
		Temperature float64 `json:"temperature_2m"`
		Time        string  `json:"time"`
		Code        int     `json:"weather_code"`
	} `json:"current"`
}

func main() {
	locationPtr := flag.String("location", "Berlin", "Location provided to get the coordinates")
	flag.Parse()

	location := *locationPtr
	if location == "" {
		fmt.Println("ERROR: Please enter a valid location with the -location Flag.")
		os.Exit(1)
	}

	encodedLocation := url.QueryEscape(location)
	geoAPIURL := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=de&format=json",
		encodedLocation,
	)
	fmt.Printf("1/2: Search for coordinates'%s'...\n", location)

	geoResponse, err := http.Get(geoAPIURL)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer geoResponse.Body.Close()

	if geoResponse.StatusCode != http.StatusOK {
		log.Fatalf("Api call exited with the error code: %d", geoResponse.StatusCode)
	}

	geoBody, err := io.ReadAll(geoResponse.Body)
	if err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	var geoData GeocodingResponse
	if err := json.Unmarshal(geoBody, &geoData); err != nil {
		log.Fatalf("Error parsing the json: %v", err)
	}

	if len(geoData.Results) == 0 {
		log.Fatalf("No coordinates found for the location'%s'.", location)
	}

	// Koordinaten des ersten Ergebnisses extrahieren
	lat := geoData.Results[0].Latitude
	lon := geoData.Results[0].Longitude
	foundName := geoData.Results[0].Name
	fmt.Printf("-> Found: %s, Lat=%.2f, Lon=%.2f\n", foundName, lat, lon)

	// cunstructing the url with the input
	api := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=temperature_2m,weather_code&timezone=Europe%%2FBerlin&forecast_days=1",
		lat,
		lon,
	)

	fmt.Printf("Sending request: %s\n", api)
	response, err := http.Get(api)
	if err != nil {
		log.Fatalf("Error sending the request: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Api call exited with the error code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error parsing the response: %v", err)
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		log.Fatalf("Error parsing the json response: %v", err)
	}

	fmt.Printf("\n-- Weather --\n")
	fmt.Printf("Lat/Lon: %.2f / %.2f\n", weatherData.Latitude, weatherData.Longitude)
	fmt.Printf("Current Temperature: %.1f °C\n", weatherData.CurrentWeather.Temperature)
	fmt.Printf("Time: %s\n", weatherData.CurrentWeather.Time)
	fmt.Printf("Weather Code: %d\n", weatherData.CurrentWeather.Code)
	fmt.Println("-------------")
}
