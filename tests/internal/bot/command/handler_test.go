package command_test

import (
	"context"
	"testing"

	"command-bot/internal/bot/command"
	pkgcommand "command-bot/pkg/command"
)

// MockCommand - это простая реализация интерфейса Command для тестирования
type MockCommand struct {
	name        string
	aliases     []string
	description string
	usage       string
	permissions []string
	executeFunc func(ctx context.Context, cmdCtx pkgcommand.CommandContext) (string, error)
}

func (c *MockCommand) Name() string                  { return c.name }
func (c *MockCommand) Aliases() []string             { return c.aliases }
func (c *MockCommand) Description() string           { return c.description }
func (c *MockCommand) Usage() string                 { return c.usage }
func (c *MockCommand) RequiredPermissions() []string { return c.permissions }

func (c *MockCommand) Execute(ctx context.Context, cmdCtx pkgcommand.CommandContext) (string, error) {
	if c.executeFunc != nil {
		return c.executeFunc(ctx, cmdCtx)
	}
	return "mock response", nil
}

func TestRegisterCommand(t *testing.T) {
	handler := command.NewHandler("/")

	// Создаем тестовую команду
	mockCmd := &MockCommand{
		name:        "test",
		aliases:     []string{"t", "tst"},
		description: "Test command",
		usage:       "test [args]",
		permissions: []string{},
	}

	// Регистрируем команду
	err := handler.RegisterCommand(mockCmd)
	if err != nil {
		t.Fatalf("Failed to register command: %v", err)
	}

	// Пытаемся получить команду по имени
	cmd, err := handler.GetCommand("test")
	if err != nil {
		t.Fatalf("Failed to get command by name: %v", err)
	}

	if cmd.Name() != "test" {
		t.Errorf("Expected command name 'test', got '%s'", cmd.Name())
	}

	// Пытаемся получить команду по псевдониму
	cmd, err = handler.GetCommand("t")
	if err != nil {
		t.Fatalf("Failed to get command by alias: %v", err)
	}

	if cmd.Name() != "test" {
		t.Errorf("Expected command name 'test', got '%s'", cmd.Name())
	}
}

func TestUnregisterCommand(t *testing.T) {
	handler := command.NewHandler("/")

	// Создаем тестовую команду
	mockCmd := &MockCommand{
		name:        "test",
		aliases:     []string{"t", "tst"},
		description: "Test command",
		usage:       "test [args]",
		permissions: []string{},
	}

	// Регистрируем команду
	err := handler.RegisterCommand(mockCmd)
	if err != nil {
		t.Fatalf("Failed to register command: %v", err)
	}

	// Отменяем регистрацию команды
	err = handler.UnregisterCommand("test")
	if err != nil {
		t.Fatalf("Failed to unregister command: %v", err)
	}

	// Пытаемся получить команду по имени (должно завершиться ошибкой)
	_, err = handler.GetCommand("test")
	if err == nil {
		t.Error("Expected error when getting unregistered command, got nil")
	}

	// Пытаемся получить команду по псевдониму (должно завершиться ошибкой)
	_, err = handler.GetCommand("t")
	if err == nil {
		t.Error("Expected error when getting unregistered command by alias, got nil")
	}
}

func TestParseCommand(t *testing.T) {
	handler := command.NewHandler("/")

	// Тестируем корректную команду
	cmdCtx, err := handler.ParseCommand("/test arg1 arg2", "user123", "chat456")
	if err != nil {
		t.Fatalf("Failed to parse valid command: %v", err)
	}

	if len(cmdCtx.Arguments) != 2 || cmdCtx.Arguments[0] != "arg1" || cmdCtx.Arguments[1] != "arg2" {
		t.Errorf("Incorrect arguments parsed: %v", cmdCtx.Arguments)
	}

	if cmdCtx.UserID != "user123" || cmdCtx.ChatID != "chat456" {
		t.Errorf("Incorrect user/chat IDs: %s/%s", cmdCtx.UserID, cmdCtx.ChatID)
	}

	// Тестируем некорректную команду (без префикса)
	_, err = handler.ParseCommand("test arg1 arg2", "user123", "chat456")
	if err == nil {
		t.Error("Expected error for command without prefix, got nil")
	}
}

func TestExecuteCommand(t *testing.T) {
	handler := command.NewHandler("/")

	// Создаем тестовую команду, которая возвращает определенный ответ
	mockCmd := &MockCommand{
		name:        "test",
		aliases:     []string{"t"},
		description: "Test command",
		usage:       "test [args]",
		permissions: []string{},
		executeFunc: func(ctx context.Context, cmdCtx pkgcommand.CommandContext) (string, error) {
			return "executed with args: " + cmdCtx.Arguments[0], nil
		},
	}

	// Регистрируем команду
	err := handler.RegisterCommand(mockCmd)
	if err != nil {
		t.Fatalf("Failed to register command: %v", err)
	}

	// Разбираем команду
	cmdCtx, err := handler.ParseCommand("/test hello", "user123", "chat456")
	if err != nil {
		t.Fatalf("Failed to parse command: %v", err)
	}

	// Выполняем команду
	response, err := handler.ExecuteCommand(context.Background(), cmdCtx)
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	expectedResponse := "executed with args: hello"
	if response != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, response)
	}
}