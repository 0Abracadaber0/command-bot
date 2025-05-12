package commands

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"command-bot/pkg/command"
)

// RandomCommand генерирует случайное число
type RandomCommand struct {
	rng *rand.Rand
}

// NewRandomCommand создает новую команду random
func NewRandomCommand() *RandomCommand {
	source := rand.NewSource(time.Now().UnixNano())
	return &RandomCommand{
		rng: rand.New(source),
	}
}

// Name возвращает основное имя команды
func (c *RandomCommand) Name() string {
	return "random"
}

// Aliases возвращает альтернативные имена для команды
func (c *RandomCommand) Aliases() []string {
	return []string{"rand", "roll"}
}

// Description возвращает краткое описание того, что делает команда
func (c *RandomCommand) Description() string {
	return "Generates a random number within a specified range"
}

// Usage возвращает строку, показывающую, как использовать команду
func (c *RandomCommand) Usage() string {
	return "random [max] or random [min] [max]"
}

// RequiredPermissions возвращает список разрешений, необходимых для выполнения этой команды
func (c *RandomCommand) RequiredPermissions() []string {
	return []string{} // Специальные разрешения не требуются
}

// Execute выполняет команду с заданным контекстом и возвращает ответ
func (c *RandomCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	var min, max int
	var err error

	switch len(cmdCtx.Arguments) {
	case 0:
		// По умолчанию от 1 до 100
		min, max = 1, 100
	case 1:
		// Только максимальное значение указано
		max, err = strconv.Atoi(cmdCtx.Arguments[0])
		if err != nil {
			return "", fmt.Errorf("invalid maximum value: %w", err)
		}
		min = 1
	case 2:
		// Указаны и минимальное, и максимальное значения
		min, err = strconv.Atoi(cmdCtx.Arguments[0])
		if err != nil {
			return "", fmt.Errorf("invalid minimum value: %w", err)
		}
		max, err = strconv.Atoi(cmdCtx.Arguments[1])
		if err != nil {
			return "", fmt.Errorf("invalid maximum value: %w", err)
		}
	default:
		return "", fmt.Errorf("too many arguments: expected 0-2, got %d", len(cmdCtx.Arguments))
	}

	if min >= max {
		return "", fmt.Errorf("minimum value must be less than maximum value")
	}

	randomNum := c.rng.Intn(max-min+1) + min
	return fmt.Sprintf("Random number between %d and %d: %d", min, max, randomNum), nil
}
