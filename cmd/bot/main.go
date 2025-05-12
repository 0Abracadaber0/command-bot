package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"command-bot/internal/bot/command"
	"command-bot/internal/bot/command/commands"
)

func main() {
	handler := command.NewHandler("/")

	pingCmd := commands.NewPingCommand()
	echoCmd := commands.NewEchoCommand()

	if err := handler.RegisterCommand(pingCmd); err != nil {
		log.Fatalf("Failed to register ping command: %v", err)
	}

	if err := handler.RegisterCommand(echoCmd); err != nil {
		log.Fatalf("Failed to register echo command: %v", err)
	}

	helpCmd := commands.NewHelpCommand(handler)
	if err := handler.RegisterCommand(helpCmd); err != nil {
		log.Fatalf("Failed to register help command: %v", err)
	}

	fmt.Println("Command Bot started. Type '/help' for available commands. Type 'exit' to quit.")

	ctx := context.Background()

	for {
		fmt.Print("> ")
		var input string
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Goodbye!")
			os.Exit(0)
		}

		cmdCtx, err := handler.ParseCommand(input, "user123", "chat456")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		response, err := handler.ExecuteCommand(ctx, cmdCtx)
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			continue
		}

		fmt.Println(response)
	}
}
