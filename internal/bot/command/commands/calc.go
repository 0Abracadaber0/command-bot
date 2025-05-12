// Пакет commands предоставляет реализации различных команд бота.
package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"command-bot/pkg/command"
)

// CalcCommand выполняет простые арифметические операции
type CalcCommand struct{}

// NewCalcCommand создает новую команду calc
func NewCalcCommand() *CalcCommand {
	return &CalcCommand{}
}

// Name возвращает основное имя команды
func (c *CalcCommand) Name() string {
	return "calc"
}

// Aliases возвращает альтернативные имена для команды
func (c *CalcCommand) Aliases() []string {
	return []string{"calculate", "math"}
}

// Description возвращает краткое описание того, что делает команда
func (c *CalcCommand) Description() string {
	return "Performs basic arithmetic calculations"
}

// Usage возвращает строку, показывающую, как использовать команду
func (c *CalcCommand) Usage() string {
	return "calc <number> <operation> <number>"
}

// RequiredPermissions возвращает список разрешений, необходимых для выполнения этой команды
func (c *CalcCommand) RequiredPermissions() []string {
	return []string{} // Специальные разрешения не требуются
}

// Execute выполняет команду с заданным контекстом и возвращает ответ
func (c *CalcCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	if len(cmdCtx.Arguments) != 3 {
		return "", fmt.Errorf("invalid format: expected 3 arguments (number operation number), got %d", len(cmdCtx.Arguments))
	}

	// Парсим первое число
	num1, err := strconv.ParseFloat(cmdCtx.Arguments[0], 64)
	if err != nil {
		return "", fmt.Errorf("invalid first number: %w", err)
	}

	// Получаем операцию
	operation := strings.ToLower(cmdCtx.Arguments[1])

	// Парсим второе число
	num2, err := strconv.ParseFloat(cmdCtx.Arguments[2], 64)
	if err != nil {
		return "", fmt.Errorf("invalid second number: %w", err)
	}

	var result float64
	var operationSymbol string

	// Выполняем операцию
	switch operation {
	case "add", "+":
		result = num1 + num2
		operationSymbol = "+"
	case "subtract", "sub", "-":
		result = num1 - num2
		operationSymbol = "-"
	case "multiply", "mul", "*", "x":
		result = num1 * num2
		operationSymbol = "×"
	case "divide", "div", "/":
		if num2 == 0 {
			return "", fmt.Errorf("division by zero is not allowed")
		}
		result = num1 / num2
		operationSymbol = "÷"
	default:
		return "", fmt.Errorf("unknown operation: %s (supported: add/+, subtract/-, multiply/*, divide/)", operation)
	}

	// Форматируем результат, убирая десятичные нули, если результат целый
	var resultStr string
	if result == float64(int(result)) {
		resultStr = fmt.Sprintf("%d", int(result))
	} else {
		resultStr = fmt.Sprintf("%.2f", result)
	}

	return fmt.Sprintf("%.2f %s %.2f = %s", num1, operationSymbol, num2, resultStr), nil
}
