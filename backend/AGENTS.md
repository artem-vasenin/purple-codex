# Backend Guidelines

## Project Structure & Module Organization

The backend is a Go service. Executables belong under `cmd/`; reusable application code belongs under `internal/`. Keep domain packages isolated from transport and infrastructure code as they are added.

## Build, Test, and Development Commands

Run commands from the `backend/` directory:

```sh
go run ./cmd/server   # Start the Go scaffold
go test ./...         # Run all Go tests
go build ./...        # Compile all Go packages
```

## Coding Style & Naming Conventions

Go code must be formatted with `gofmt`. Use idiomatic mixedCaps names and short, focused packages.

## Testing Guidelines

Use the standard library `testing` package. Test files should be named `*_test.go`, with test functions such as `TestServerStarts`.
