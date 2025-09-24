# Go Code Style and Architectural Guide

This document outlines the Go-specific code style, conventions, and architectural patterns for the `backlog` repository. Adhering to these guidelines is crucial for maintaining a clean, idiomatic, and maintainable codebase.

## 1. Project Layout

The project follows the standard Go project layout recommendations.

-   **`internal/`**: All core application code is located here. This code is not meant to be imported by other projects.
    -   `internal/cmd`: Defines the CLI commands using the **Cobra** library. Each command and its related flags are in a separate file (e.g., `task_create.go`).
    -   `internal/core`: Contains the core business logic and data structures (e.g., `Task`, `TaskStore`). This package is the heart of the application and is designed to be decoupled from the CLI.
    -   `internal/commit`: Handles Git-related operations, such as automatic commits.
    -   `internal/mcp`: Implements the Model-Context-Protocol (MCP) server for AI agent integration.
    -   `internal/logging`: Provides a structured logging setup for the application.
    -   `internal/paths`: Provides utilities for resolving repository paths.
-   **`pkg/`**: (Not currently used) Would contain library code intended for external use.
-   **`main.go`**: The main application entry point. It is kept minimal, with its primary role being to initialize and execute the root Cobra command.

## 2. Go Language & Idioms

-   **Formatting**: All code **must** be formatted with `gofmt`. The CI pipeline will fail if code is not formatted correctly.
-   **Linting**: We use `golangci-lint` with a strict configuration to enforce idiomatic Go. Run `make lint` locally before pushing changes.
-   **Naming Conventions**:
    -   Package names are short, concise, and all lowercase (e.g., `core`, `commit`).
    -   Public symbols (variables, functions, types) are `PascalCase`.
    -   Private symbols are `camelCase`.
    -   Acronyms like `ID` and `API` should be consistently cased (e.g., `taskID`, not `taskId`).
-   **Interfaces**:
    -   Interfaces are defined by the consumer. For example, if a function needs a `Reader`, it should accept `io.Reader`, not a concrete type.
    -   Keep interfaces small and focused on a single behavior (e.g., `io.Reader`, `fmt.Stringer`). The `TaskStore` interface in `internal/core/store.go` is a good example.
-   **Structs**:
    -   Structs should be initialized with explicit field names where possible (e.g., `core.Task{Title: "New Task"}`).
    -   Group related fields together. Add comments to explain complex or non-obvious fields.
-   **Pointers vs. Values**:
    -   Use pointers for large structs or when a method needs to modify the receiver.
    -   Use values for small, immutable structs or built-in types.
    -   In this project, `*Task` is used frequently because tasks are often modified and passed around.

## 3. Error Handling

-   **Always Check Errors**: Never ignore an error with the blank identifier (`_`). The only exception is for a `Close()` call on a read-only resource where failure has no consequence.
-   **Error Wrapping**: Errors that cross package boundaries **must** be wrapped to provide context. Use `fmt.Errorf("operation failed: %w", err)`. This creates a chain of errors that can be inspected for debugging.
-   **Custom Error Types**: For specific, expected errors (e.g., "task not found"), use custom error variables like `core.ErrNotFound`. This allows callers to check for specific error conditions using `errors.Is()`.
-   **No Panics**: The application **must not** use `panic` for recoverable errors. Panics are reserved for unrecoverable, programmer-level mistakes that indicate a bug.

## 4. Logging

-   **Structured Logging**: We use a structured logger (e.g., `zerolog`) for all application output. This allows for easier parsing, filtering, and analysis of logs.
-   **Log Levels**: Use appropriate log levels:
    -   `logging.Error`: For failures that prevent a feature from working.
    -   `logging.Warn`: For potential issues that do not cause a failure.
    -   `logging.Info`: For general, informative messages (e.g., "task created successfully").
    -   `logging.Debug`: For verbose, development-only messages.

## 5. Concurrency

-   (Not heavily used yet) When introducing concurrency, prefer channels for communication and synchronization over explicit locks where possible.
-   Ensure all goroutines are properly managed and have a clear exit path to prevent leaks.

## 6. Testing

-   **Test Location**: Tests for `foo.go` are located in `foo_test.go` within the same package.
-   **Test Tables**: For testing multiple inputs and outputs, use table-driven tests.
-   **Filesystem Abstraction**: We use the **`afero`** library to abstract the filesystem.
    -   In production, `afero.NewOsFs()` is used.
    -   In tests, `afero.NewMemMapFs()` provides a fast, in-memory filesystem, making tests hermetic and reliable.
-   **Dependency Injection**: Core components like `FileTaskStore` are designed for dependency injection. They accept dependencies (like `afero.Fs`) via their constructor (`NewFileTaskStore`). This is critical for decoupling and testability.

## 7. CLI Implementation

-   **Library**: The CLI is built using the `cobra` library.
-   **Command Documentation**: All user-facing commands must include a `Short` description, a `Long` description, and comprehensive `Example` usage strings.

## 8. Git Integration

-   **Automatic Commits**: Task operations (create, edit, archive) trigger automatic Git commits to maintain a history of changes.
-   **Clean Worktree**: The auto-commit feature will not run if the Git worktree is dirty. This prevents accidental inclusion of unrelated changes.

## 9. Code Generation

-   **`go:generate`**: We use `//go:generate` to automate tasks like generating CLI documentation. Run `make docs` to update generated content. This ensures documentation stays in sync with the code.
