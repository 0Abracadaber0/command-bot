package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"command-bot/pkg/command"
)

// WeatherCommand получает информацию о текущей погоде из Open-Meteo
// и не требует API-ключа
type WeatherCommand struct{}

// NewWeatherCommand создает новую команду weather
func NewWeatherCommand() *WeatherCommand {
	return &WeatherCommand{}
}

// Name возвращает основное имя команды
func (c *WeatherCommand) Name() string {
	return "weather"
}

// Aliases возвращает альтернативные имена для команды
func (c *WeatherCommand) Aliases() []string {
	return []string{"forecast", "temp"}
}

// Description возвращает краткое описание того, что делает команда
func (c *WeatherCommand) Description() string {
	return "Shows current weather information for a specified location using Open-Meteo (no API key required)"
}

// Usage возвращает строку, показывающую, как использовать команду
func (c *WeatherCommand) Usage() string {
	return "weather <location>"
}

// RequiredPermissions возвращает список разрешений, необходимых для выполнения этой команды
func (c *WeatherCommand) RequiredPermissions() []string {
	return []string{}
}

// Execute выполняет команду с заданным контекстом и возвращает ответ
func (c *WeatherCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	if len(cmdCtx.Arguments) == 0 {
		return "", fmt.Errorf("please specify a location")
	}

	// Собираем название локации
	location := strings.Join(cmdCtx.Arguments, " ")
	// Кодируем параметр
	q := url.QueryEscape(location)

	// Шаг 1: геокодирование через Open-Meteo Geocoding API
	geoURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", q)
	resp, err := http.Get(geoURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch geocoding data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	var geoData struct {
		Results []struct {
			Name      string  `json:"name"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		return "", fmt.Errorf("failed to parse geocoding response: %v", err)
	}
	if len(geoData.Results) == 0 {
		return "", fmt.Errorf("location not found: %s", location)
	}
	city := geoData.Results[0].Name
	country := geoData.Results[0].Country
	lat := geoData.Results[0].Latitude
	lon := geoData.Results[0].Longitude

	// Шаг 2: запрос текущей погоды
	weatherURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current_weather=true&timezone=Europe/Warsaw",
		lat, lon,
	)
	resp2, err := http.Get(weatherURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return "", fmt.Errorf("weather API returned status %d", resp2.StatusCode)
	}

	var weatherData struct {
		CurrentWeather struct {
			Temperature   float64 `json:"temperature"`
			Windspeed     float64 `json:"windspeed"`
			Winddirection int     `json:"winddirection"`
			Weathercode   int     `json:"weathercode"`
			Time          string  `json:"time"`
		} `json:"current_weather"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&weatherData); err != nil {
		return "", fmt.Errorf("failed to parse weather response: %v", err)
	}

	cw := weatherData.CurrentWeather

	return fmt.Sprintf(
		"Weather for %s, %s:\n"+
			"Time: %s\n"+
			"Temperature: %.1f°C\n"+
			"Wind Speed: %.1f m/s (direction %d°)\n"+
			"Weather Code: %d",
		city, country,
		cw.Time,
		cw.Temperature,
		cw.Windspeed,
		cw.Winddirection,
		cw.Weathercode,
	), nil
}
