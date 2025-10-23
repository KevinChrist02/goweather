package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
	api := "https://api.open-meteo.com/v1/forecast?latitude=49.3064&longitude=8.4503&current=temperature_2m,weather_code&timezone=Europe%2FBerlin&forecast_days=1"

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
		log.Fatalf("Error reading the response: %v", err)
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
