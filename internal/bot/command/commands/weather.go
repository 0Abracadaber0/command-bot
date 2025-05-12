package commands

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"command-bot/pkg/command"
)

// WeatherCommand симулирует предоставление информации о погоде
type WeatherCommand struct {
	rng *rand.Rand
}

// NewWeatherCommand создает новую команду weather
func NewWeatherCommand() *WeatherCommand {
	source := rand.NewSource(time.Now().UnixNano())
	return &WeatherCommand{
		rng: rand.New(source),
	}
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
	return "Shows simulated weather information for a specified location"
}

// Usage возвращает строку, показывающую, как использовать команду
func (c *WeatherCommand) Usage() string {
	return "weather <location>"
}

// RequiredPermissions возвращает список разрешений, необходимых для выполнения этой команды
func (c *WeatherCommand) RequiredPermissions() []string {
	return []string{} // Специальные разрешения не требуются
}

// Execute выполняет команду с заданным контекстом и возвращает ответ
func (c *WeatherCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	if len(cmdCtx.Arguments) == 0 {
		return "", fmt.Errorf("please specify a location")
	}

	location := strings.Join(cmdCtx.Arguments, " ")

	// Симулируем различные погодные условия
	conditions := []string{
		"Sunny", "Cloudy", "Partly Cloudy", "Rainy", "Stormy",
		"Snowy", "Foggy", "Windy", "Clear", "Overcast",
	}

	// Генерируем случайную температуру от -20 до 40 градусов Цельсия
	temperature := c.rng.Intn(61) - 20

	// Генерируем случайную влажность от 0 до 100%
	humidity := c.rng.Intn(101)

	// Генерируем случайную скорость ветра от 0 до 30 км/ч
	windSpeed := c.rng.Intn(31)

	// Выбираем случайное погодное условие
	condition := conditions[c.rng.Intn(len(conditions))]

	return fmt.Sprintf("Weather for %s:\n"+
		"Condition: %s\n"+
		"Temperature: %d°C\n"+
		"Humidity: %d%%\n"+
		"Wind Speed: %d km/h",
		location, condition, temperature, humidity, windSpeed), nil
}
