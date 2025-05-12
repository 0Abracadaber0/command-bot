// Пакет commands предоставляет реализации различных команд бота.
package commands

import (
	"context"
	"fmt"
	"time"

	"command-bot/pkg/command"
)

// PingCommand отвечает сообщением pong и информацией о задержке
type PingCommand struct{}

// NewPingCommand создает новую команду ping
func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

// Name возвращает основное имя команды
func (c *PingCommand) Name() string {
	return "ping"
}

// Aliases возвращает альтернативные имена для команды
func (c *PingCommand) Aliases() []string {
	return []string{"latency", "p"}
}

// Description возвращает краткое описание того, что делает команда
func (c *PingCommand) Description() string {
	return "Checks if the bot is responsive and shows latency"
}

// Usage возвращает строку, показывающую, как использовать команду
func (c *PingCommand) Usage() string {
	return "ping"
}

// RequiredPermissions возвращает список разрешений, необходимых для выполнения этой команды
func (c *PingCommand) RequiredPermissions() []string {
	return []string{} // Специальные разрешения не требуются
}

// Execute выполняет команду с заданным контекстом и возвращает ответ
func (c *PingCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	// Записываем время начала
	start := time.Now()

	// Имитируем некоторое время обработки
	time.Sleep(10 * time.Millisecond)

	// Вычисляем задержку
	latency := time.Since(start)

	return fmt.Sprintf("Pong! Latency: %v", latency), nil
}
