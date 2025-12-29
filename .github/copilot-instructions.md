# Folio Project Context

## Project Overview
Folio is a CLI tool for preparing PDFs for printing as booklets or multi-up handouts. It is written in Go.

## Tech Stack
- **Language:** Go (v1.25.1)
- **CLI Framework:** [kong](https://github.com/alecthomas/kong)
- **PDF Processing:** [pdfcpu](https://github.com/pdfcpu/pdfcpu)

## Directory Structure
- `cmd/folio/`: Entry point for the application (`main.go`).
- `internal/cli/`: CLI command definitions and logic (using `kong`).
- `internal/pdf/`: PDF processing logic and wrappers around `pdfcpu`.
- `internal/utils/`: General utility functions (e.g., error handling).
- `testdata/`: Sample files for testing.

## Development Workflow

### Pre-Commit Checklist
Before committing any changes, ensure the following steps are performed to maintain code quality and consistency:

1.  **Linting & Formatting:**
    Run `golangci-lint` to check for linting errors and apply formatting.
    ```bash
    golangci-lint run
    ```
    *Note: The project uses `govet`, `staticcheck`, and `ineffassign`. Formatting is handled by `gofmt` and `goimports`.*

2.  **Static Analysis:**
    Run `go vet` to catch suspicious constructs.
    ```bash
    go vet ./...
    ```

3.  **Testing:**
    Run all unit tests.
    ```bash
    go test ./...
    ```

4.  **Dependency Management:**
    Ensure `go.mod` and `go.sum` are tidy.
    ```bash
    go mod tidy
    ```

## Coding Conventions
- **Error Handling:** Use `internal/utils/error.go` helpers (`WrapErr`, `NewErr`) for consistent error wrapping and formatting.
- **File Validation:** Ensure PDF files are validated (extension, content, EOF marker) before processing.
- **CLI Structure:** Define commands as structs in `internal/cli/` and register them in the main CLI struct.

## Common Commands
- **Run App:** `go run cmd/folio/main.go`
- **Run Tests:** `go test -v ./...`
