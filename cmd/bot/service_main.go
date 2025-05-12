package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"command-bot/internal/bot/command"
	"command-bot/internal/bot/command/commands"
)

func main() {
	logFile, err := os.OpenFile("/var/log/command-bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Command Bot service starting...")

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

	log.Println("Command Bot service started. Registered commands: ping, echo, help")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v, shutting down...", sig)
		cancel()
	}()

	type CommandRequest struct {
		Command string `json:"command"`
		UserID  string `json:"user_id,omitempty"`
		ChatID  string `json:"chat_id,omitempty"`
	}

	type CommandResponse struct {
		Response string `json:"response"`
		Error    string `json:"error,omitempty"`
	}

	http.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
			return
		}

		userID := req.UserID
		if userID == "" {
			userID = "api-user"
		}

		chatID := req.ChatID
		if chatID == "" {
			chatID = "api-chat"
		}

		log.Printf("Received command: %s from user: %s in chat: %s", req.Command, userID, chatID)

		cmdCtx, err := handler.ParseCommand(req.Command, userID, chatID)
		if err != nil {
			response := CommandResponse{
				Error: fmt.Sprintf("Error parsing command: %v", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		cmdResponse, err := handler.ExecuteCommand(ctx, cmdCtx)
		if err != nil {
			response := CommandResponse{
				Error: fmt.Sprintf("Error executing command: %v", err),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := CommandResponse{
			Response: cmdResponse,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Command Bot service is running")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Command Bot service shutting down")
}
