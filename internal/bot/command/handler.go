package command

import (
	"context"
	"strings"
	"sync"

	"command-bot/pkg/command"
)

type Handler struct {
	commands map[string]command.Command
	aliases  map[string]string
	prefix   string
	mu       sync.RWMutex
}

// NewHandler создает новый обработчик команд с заданным префиксом
func NewHandler(prefix string) *Handler {
	return &Handler{
		commands: make(map[string]command.Command),
		aliases:  make(map[string]string),
		prefix:   prefix,
	}
}

// RegisterCommand добавляет команду в обработчик
func (h *Handler) RegisterCommand(cmd command.Command) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	name := strings.ToLower(cmd.Name())

	if _, exists := h.commands[name]; exists {
		return command.ErrCommandNotFound
	}

	h.commands[name] = cmd

	for _, alias := range cmd.Aliases() {
		alias = strings.ToLower(alias)
		h.aliases[alias] = name
	}

	return nil
}

func (h *Handler) UnregisterCommand(name string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	name = strings.ToLower(name)

	cmd, exists := h.commands[name]
	if !exists {
		return command.ErrCommandNotFound
	}

	delete(h.commands, name)

	for _, alias := range cmd.Aliases() {
		alias = strings.ToLower(alias)
		delete(h.aliases, alias)
	}

	return nil
}

// GetCommand получает команду по имени или псевдониму
func (h *Handler) GetCommand(name string) (command.Command, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	name = strings.ToLower(name)

	if cmd, exists := h.commands[name]; exists {
		return cmd, nil
	}

	if primaryName, exists := h.aliases[name]; exists {
		return h.commands[primaryName], nil
	}

	return nil, command.ErrCommandNotFound
}

// ListCommands возвращает все зарегистрированные команды
func (h *Handler) ListCommands() []command.Command {
	h.mu.RLock()
	defer h.mu.RUnlock()

	commands := make([]command.Command, 0, len(h.commands))
	for _, cmd := range h.commands {
		commands = append(commands, cmd)
	}

	return commands
}

// ParseCommand разбирает необработанный ввод в CommandContext
func (h *Handler) ParseCommand(input string, userID, chatID string) (command.CommandContext, error) {
	if !strings.HasPrefix(input, h.prefix) {
		return command.CommandContext{}, command.ErrCommandNotFound
	}

	input = strings.TrimPrefix(input, h.prefix)

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return command.CommandContext{}, command.ErrCommandNotFound
	}

	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	return command.CommandContext{
		UserID:    userID,
		ChatID:    chatID,
		Arguments: args,
		RawInput:  input,
		Metadata:  make(map[string]interface{}),
	}, nil
}

// ExecuteCommand обрабатывает и выполняет команду
func (h *Handler) ExecuteCommand(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	parts := strings.Fields(cmdCtx.RawInput)
	if len(parts) == 0 {
		return "", command.ErrCommandNotFound
	}

	cmdName := parts[0]

	cmd, err := h.GetCommand(cmdName)
	if err != nil {
		return "", err
	}

	return cmd.Execute(ctx, cmdCtx)
}
