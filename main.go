package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (ApiConfigData, error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return ApiConfigData{}, err
	}
	var ConfigData ApiConfigData

	err = json.Unmarshal(data, &ConfigData)
	if err != nil {
		fmt.Println(err)
		return ApiConfigData{}, err
	}
	return ConfigData, nil
}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Go and Omar! \n"))
}
func query(city string) (weatherData, error) {
	ApiConfigData, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + ApiConfigData.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var data weatherData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return weatherData{}, err
	}
	return data, nil
}
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome! Try /hello or /weather/{city}"))
	})
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		City := strings.SplitN(r.URL.Path, "/", 3)[2]
		CityData, err := query(City)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CityData)
	})
	http.ListenAndServe(":8080", nil)
}
