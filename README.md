# Command Bot

A chat bot with command support module following the standard Go project layout.

## Project Structure

```
.
├── .github/workflows # GitHub Actions workflow configurations
├── api/proto         # Protocol buffer definitions
├── cmd/bot           # Main application entry points
├── docs              # Documentation
├── examples          # Example code
├── internal          # Private application and library code
│   └── bot
│       └── command   # Internal command handling logic
├── pkg               # Library code that can be used by external applications
│   └── command       # Public command handling interfaces and utilities
└── tests             # Test files mirroring the package structure
    └── internal
        └── bot
            └── command # Tests for internal command handling logic
```

## Overview

This project implements a chat bot with a focus on command support functionality. It provides a framework for registering, discovering, and executing commands within a chat context.

## Features

- Command registration and discovery
- Command execution with parameter parsing
- Help system for available commands
- Extensible command framework

## Usage

To use this bot, run:

```bash
go run cmd/bot/main.go
```

## Development

To add a new command:

1. Implement the command interface in `pkg/command`
2. Register the command in the command registry
3. The command will be automatically available to users

## Testing

Tests are located in the `tests` directory, mirroring the package structure of the code being tested.

To run tests:

```bash
go test ./tests/...
```

To run tests with coverage:

```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## CI/CD

This project uses GitHub Actions for continuous integration. The workflow is defined in `.github/workflows/go-test.yml` and includes:

- Running tests on push to main branch and pull requests
- Generating and uploading test coverage reports

## Running as a Linux Service

This project includes two versions of the application:
- `main.go`: Interactive version for terminal use
- `service_main.go`: Non-interactive version designed to run as a service

### Using the Service Version

To run the Command Bot as a systemd service:

1. Build the service version of the application:
   ```bash
   go build -o command-bot-service cmd/bot/service_main.go
   ```

2. Create a directory for the application:
   ```bash
   sudo mkdir -p /opt/command-bot
   ```

3. Copy the executable and any necessary files:
   ```bash
   sudo cp command-bot-service /opt/command-bot/
   ```

4. Copy the service file to the systemd directory:
   ```bash
   sudo cp command-bot.service /etc/systemd/system/
   ```

5. Reload systemd to recognize the new service:
   ```bash
   sudo systemctl daemon-reload
   ```

6. Enable the service to start on boot:
   ```bash
   sudo systemctl enable command-bot.service
   ```

7. Start the service:
   ```bash
   sudo systemctl start command-bot.service
   ```

8. Check the service status:
   ```bash
   sudo systemctl status command-bot.service
   ```

9. View logs:
   ```bash
   sudo journalctl -u command-bot.service
   # or check the application log file
   sudo cat /var/log/command-bot.log
   ```

### About the Service Version

The service version (`service_main.go`):
- Logs to `/var/log/command-bot.log` instead of stdout
- Handles OS signals for graceful shutdown
- Provides an HTTP API for sending commands to the bot
- Is suitable for running as a background service

### Sending Commands to the Service

When running as a service, the Command Bot exposes an HTTP API on port 8080 that allows you to send commands to it. Here's how to use it:

#### API Endpoints

1. **Command Endpoint**: `/command`
   - Method: POST
   - Content-Type: application/json
   - Request Body:
     ```json
     {
       "command": "/your-command-here",
       "user_id": "optional-user-id",
       "chat_id": "optional-chat-id"
     }
     ```
   - Response:
     ```json
     {
       "response": "Command response text",
       "error": "Error message (if any)"
     }
     ```

2. **Health Check**: `/health`
   - Method: GET
   - Response: Plain text indicating the service is running

#### Example Usage

To send a command to the bot, you can use curl:

```bash
# Send a ping command
curl -X POST http://localhost:8080/command \
  -H "Content-Type: application/json" \
  -d '{"command": "/ping"}'

# Send an echo command
curl -X POST http://localhost:8080/command \
  -H "Content-Type: application/json" \
  -d '{"command": "/echo Hello, world!"}'

# Check if the service is running
curl http://localhost:8080/health
```

You can also use any HTTP client library in your application to send commands to the bot.

#### Security Considerations

The HTTP API does not include authentication by default. If you're deploying this in a production environment, consider:

1. Adding authentication to the API endpoints
2. Using HTTPS instead of HTTP
3. Restricting access to the API using a firewall or reverse proxy

You can modify `service_main.go` to add these security features or to integrate with other input sources like message queues or other communication channels.

## License

MIT
