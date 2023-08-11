package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var weatherAPIURL = "http://api.openweathermap.org/data/2.5/weather"

// getWeatherByLocation делает запрос к API погоды и возвращает информацию о погоде.
func GetWeatherByLocation(lat float64, lon float64, apiKey string) (string, error) {
	// Формирование URL для запроса к API
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&lang=ru&units=metric", weatherAPIURL, lat, lon, apiKey)

	// Выполнение GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make request to weather API: %v", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Парсинг ответа
	var weatherResponse struct {
		Name string `json:"name"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	}

	log.Printf("Weather API response: %s", body)

	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse weather API response: %v", err)
	}

	if len(weatherResponse.Weather) == 0 {
		return "", fmt.Errorf("no weather data received for this location")
	}
	// Формирование ответа для пользователя
	result := fmt.Sprintf("Город: %s\nОписание погоды: %s\nТемпература: %.2f°C",
		weatherResponse.Name,
		weatherResponse.Weather[0].Description,
		weatherResponse.Main.Temp)

	return result, nil
}

// getWeatherByCityName делает запрос к API погоды и возвращает информацию о погоде для указанного города.
func GetWeatherByCityName(client HTTPClient, cityName string, apiKey string) (string, error) {
	// Формирование URL для запроса к API
	url := fmt.Sprintf("%s?q=%s&appid=%s&lang=ru&units=metric", weatherAPIURL, cityName, apiKey)

	// Выполнение GET-запроса
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request to weather API: %v", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Парсинг ответа
	var weatherResponse struct {
		Name string `json:"name"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	}

	log.Printf("Weather API response: %s", body)

	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse weather API response: %v", err)
	}

	if len(weatherResponse.Weather) == 0 {
		return "", fmt.Errorf("no weather data received for the city: %s", cityName)
	}

	// Формирование ответа для пользователя
	result := fmt.Sprintf("Город: %s\nОписание погоды: %s\nТемпература: %.2f°C",
		weatherResponse.Name,
		weatherResponse.Weather[0].Description,
		weatherResponse.Main.Temp)

	return result, nil
}
