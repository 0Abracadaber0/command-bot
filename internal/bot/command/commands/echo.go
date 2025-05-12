package commands

import (
	"context"
	"strings"

	"command-bot/pkg/command"
)

type EchoCommand struct{}

// NewEchoCommand создает новую команду echo
func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

func (c *EchoCommand) Name() string {
	return "echo"
}

func (c *EchoCommand) Aliases() []string {
	return []string{"repeat", "say"}
}

func (c *EchoCommand) Description() string {
	return "Repeats the text you provide"
}

func (c *EchoCommand) Usage() string {
	return "echo <text>"
}

func (c *EchoCommand) RequiredPermissions() []string {
	return []string{}
}

func (c *EchoCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	if len(cmdCtx.Arguments) == 0 {
		return "You didn't provide any text to echo!", nil
	}

	return strings.Join(cmdCtx.Arguments, " "), nil
}
