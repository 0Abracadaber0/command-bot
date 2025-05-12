// Пакет commands предоставляет реализации различных команд бота.
package commands

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"command-bot/pkg/command"
)

// HelpCommand предоставляет информацию о доступных командах
type HelpCommand struct {
	handler command.CommandHandler
}

// NewHelpCommand создает новую команду help
func NewHelpCommand(handler command.CommandHandler) *HelpCommand {
	return &HelpCommand{
		handler: handler,
	}
}

// Name возвращает основное имя команды
func (c *HelpCommand) Name() string {
	return "help"
}

// Aliases возвращает альтернативные имена для команды
func (c *HelpCommand) Aliases() []string {
	return []string{"h", "?"}
}

// Description возвращает краткое описание того, что делает команда
func (c *HelpCommand) Description() string {
	return "Displays help information about available commands"
}

// Usage возвращает строку, показывающую, как использовать команду
func (c *HelpCommand) Usage() string {
	return "help [command]"
}

// RequiredPermissions возвращает список разрешений, необходимых для выполнения этой команды
func (c *HelpCommand) RequiredPermissions() []string {
	return []string{} // Специальные разрешения не требуются
}

// Execute выполняет команду с заданным контекстом и возвращает ответ
func (c *HelpCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	// Если запрошена конкретная команда, показываем справку для этой команды
	if len(cmdCtx.Arguments) > 0 {
		cmdName := cmdCtx.Arguments[0]
		cmd, err := c.handler.GetCommand(cmdName)
		if err != nil {
			return "", fmt.Errorf("command '%s' not found: %w", cmdName, err)
		}

		return c.formatCommandHelp(cmd), nil
	}

	// В противном случае, выводим список всех доступных команд
	return c.listAllCommands(), nil
}

// formatCommandHelp форматирует подробную справку для конкретной команды
func (c *HelpCommand) formatCommandHelp(cmd command.Command) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Command: %s\n", cmd.Name()))

	aliases := cmd.Aliases()
	if len(aliases) > 0 {
		sb.WriteString(fmt.Sprintf("Aliases: %s\n", strings.Join(aliases, ", ")))
	}

	sb.WriteString(fmt.Sprintf("Description: %s\n", cmd.Description()))
	sb.WriteString(fmt.Sprintf("Usage: %s\n", cmd.Usage()))

	permissions := cmd.RequiredPermissions()
	if len(permissions) > 0 {
		sb.WriteString(fmt.Sprintf("Required Permissions: %s\n", strings.Join(permissions, ", ")))
	} else {
		sb.WriteString("Required Permissions: None\n")
	}

	return sb.String()
}

// listAllCommands форматирует список всех доступных команд
func (c *HelpCommand) listAllCommands() string {
	var sb strings.Builder

	sb.WriteString("Available Commands:\n")

	// Получаем все команды и сортируем их по имени
	commands := c.handler.ListCommands()
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name() < commands[j].Name()
	})

	// Форматируем каждую команду
	for _, cmd := range commands {
		sb.WriteString(fmt.Sprintf("  %-15s %s\n", cmd.Name(), cmd.Description()))
	}

	sb.WriteString("\nType 'help <command>' for more information about a specific command.")

	return sb.String()
}
