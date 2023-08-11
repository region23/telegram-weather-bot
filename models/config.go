package models

type Config struct {
    Token         string `yaml:"token"`
    WeatherAPIKey string `yaml:"weather_api_key"`
}
