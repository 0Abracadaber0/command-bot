package command

import (
	"context"
	"errors"
)

// Общие ошибки
var (
	ErrCommandNotFound        = errors.New("command not found")
	ErrInvalidArguments       = errors.New("invalid command arguments")
	ErrCommandExecutionFailed = errors.New("command execution failed")
	ErrPermissionDenied       = errors.New("permission denied to execute command")
)

type CommandContext struct {
	UserID    string
	ChatID    string
	Arguments []string
	RawInput  string
	Metadata  map[string]interface{}
}

type Command interface {
	Name() string
	Aliases() []string
	Description() string
	Usage() string

	Execute(ctx context.Context, cmdCtx CommandContext) (string, error)
	RequiredPermissions() []string
}

type CommandHandler interface {
	RegisterCommand(cmd Command) error
	UnregisterCommand(name string) error
	GetCommand(name string) (Command, error)
	ListCommands() []Command
	ExecuteCommand(ctx context.Context, cmdCtx CommandContext) (string, error)
	ParseCommand(input string, userID, chatID string) (CommandContext, error)
}
