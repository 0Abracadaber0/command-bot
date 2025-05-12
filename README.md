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

## License

MIT
