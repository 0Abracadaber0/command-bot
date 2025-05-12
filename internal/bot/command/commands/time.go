package commands

import (
	"context"
	"fmt"
	"time"

	"command-bot/pkg/command"
)

// TimeCommand отображает текущее время
type TimeCommand struct{}

// NewTimeCommand создает новую команду time
func NewTimeCommand() *TimeCommand {
	return &TimeCommand{}
}

// Name возвращает основное имя команды
func (c *TimeCommand) Name() string {
	return "time"
}

// Aliases возвращает альтернативные имена для команды
func (c *TimeCommand) Aliases() []string {
	return []string{"now", "date"}
}

// Description возвращает краткое описание того, что делает команда
func (c *TimeCommand) Description() string {
	return "Shows the current date and time"
}

// Usage возвращает строку, показывающую, как использовать команду
func (c *TimeCommand) Usage() string {
	return "time"
}

// RequiredPermissions возвращает список разрешений, необходимых для выполнения этой команды
func (c *TimeCommand) RequiredPermissions() []string {
	return []string{} // Специальные разрешения не требуются
}

// Execute выполняет команду с заданным контекстом и возвращает ответ
func (c *TimeCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	now := time.Now()
	return fmt.Sprintf("Current time: %s", now.Format("2006-01-02 15:04:05")), nil
}
