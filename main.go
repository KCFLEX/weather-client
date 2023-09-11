package main

import (
	"encoding/json" // to parse json
	"io"            // allows basic basic read and write commands

	"net/http" //allows us to make get requests
	"os"
	"strings" // to work with strings
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

// format in which the weather data will be shown
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

// the function to get apikey from file
func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.Open(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	defer bytes.Close()

	// read the entire file content
	data, err := io.ReadAll(bytes)
	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData
	err = json.Unmarshal(data, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, err
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("nice to meet you from go\n"))
}

// query func
func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.Split(r.URL.Path, "/")[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})
	http.ListenAndServe(":8080", nil)
}
