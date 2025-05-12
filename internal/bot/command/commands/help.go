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

	// Категории команд для организации вывода справки
	categories map[string][]string
}

// NewHelpCommand создает новую команду help
func NewHelpCommand(handler command.CommandHandler) *HelpCommand {
	// Определяем категории команд
	categories := map[string][]string{
		"Utility":     {"help", "echo"},
		"Information": {"ping", "time", "weather"},
		"Fun":         {"random", "quote", "calc"},
	}

	return &HelpCommand{
		handler:    handler,
		categories: categories,
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

	// Заголовок с именем команды
	sb.WriteString(fmt.Sprintf("=== Command: %s ===\n\n", cmd.Name()))

	// Псевдонимы команды
	aliases := cmd.Aliases()
	if len(aliases) > 0 {
		sb.WriteString(fmt.Sprintf("Aliases: %s\n\n", strings.Join(aliases, ", ")))
	}

	// Описание команды
	sb.WriteString(fmt.Sprintf("Description: %s\n\n", cmd.Description()))

	// Использование команды
	sb.WriteString(fmt.Sprintf("Usage: %s\n\n", cmd.Usage()))

	// Примеры использования команды
	sb.WriteString("Examples:\n")

	// Добавляем примеры в зависимости от типа команды
	switch cmd.Name() {
	case "help":
		sb.WriteString("  /help           - Shows list of all available commands\n")
		sb.WriteString("  /help ping      - Shows detailed help for the ping command\n")
	case "echo":
		sb.WriteString("  /echo Hello World  - Bot responds with 'Hello World'\n")
		sb.WriteString("  /echo I am a bot   - Bot responds with 'I am a bot'\n")
	case "ping":
		sb.WriteString("  /ping           - Bot responds with 'Pong!' and latency information\n")
	case "time":
		sb.WriteString("  /time           - Bot responds with the current date and time\n")
	case "random":
		sb.WriteString("  /random         - Generates a random number between 1 and 100\n")
		sb.WriteString("  /random 50      - Generates a random number between 1 and 50\n")
		sb.WriteString("  /random 10 20   - Generates a random number between 10 and 20\n")
	case "weather":
		sb.WriteString("  /weather Moscow    - Shows simulated weather for Moscow\n")
		sb.WriteString("  /weather New York  - Shows simulated weather for New York\n")
	case "calc":
		sb.WriteString("  /calc 5 + 3     - Calculates 5 + 3 = 8\n")
		sb.WriteString("  /calc 10 - 4    - Calculates 10 - 4 = 6\n")
		sb.WriteString("  /calc 7 * 2     - Calculates 7 * 2 = 14\n")
		sb.WriteString("  /calc 20 / 5    - Calculates 20 / 5 = 4\n")
	case "quote":
		sb.WriteString("  /quote          - Shows a random inspirational quote\n")
	default:
		sb.WriteString("  No examples available for this command.\n")
	}

	sb.WriteString("\n")

	// Требуемые разрешения
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

	sb.WriteString("=== Available Commands ===\n\n")

	// Получаем все команды
	allCommands := c.handler.ListCommands()

	// Создаем карту для быстрого доступа к командам по имени
	cmdMap := make(map[string]command.Command)
	for _, cmd := range allCommands {
		cmdMap[cmd.Name()] = cmd
	}

	// Отслеживаем команды, которые уже были отображены в категориях
	displayedCommands := make(map[string]bool)

	// Отображаем команды по категориям
	for category, cmdNames := range c.categories {
		sb.WriteString(fmt.Sprintf("== %s ==\n", category))

		// Сортируем имена команд в категории
		sort.Strings(cmdNames)

		for _, cmdName := range cmdNames {
			if cmd, exists := cmdMap[cmdName]; exists {
				sb.WriteString(fmt.Sprintf("  %-15s %s\n", cmd.Name(), cmd.Description()))
				displayedCommands[cmdName] = true
			}
		}

		sb.WriteString("\n")
	}

	// Проверяем, есть ли команды, которые не попали ни в одну категорию
	var uncategorizedCmds []command.Command
	for _, cmd := range allCommands {
		if !displayedCommands[cmd.Name()] {
			uncategorizedCmds = append(uncategorizedCmds, cmd)
		}
	}

	// Если есть некатегоризированные команды, отображаем их отдельно
	if len(uncategorizedCmds) > 0 {
		sb.WriteString("== Other ==\n")

		// Сортируем некатегоризированные команды по имени
		sort.Slice(uncategorizedCmds, func(i, j int) bool {
			return uncategorizedCmds[i].Name() < uncategorizedCmds[j].Name()
		})

		for _, cmd := range uncategorizedCmds {
			sb.WriteString(fmt.Sprintf("  %-15s %s\n", cmd.Name(), cmd.Description()))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("Type 'help <command>' for more information about a specific command.")

	return sb.String()
}
