
# Comprehensive Guide to Idiomatic Go for LLMs

## Introduction

This document provides a comprehensive guide for Large Language Models (LLMs) on how to write idiomatic, maintainable, and testable Go code. It is a synthesis of best practices from several authoritative sources in the Go community. By following these guidelines, you will be able to generate Go code that is not only correct but also aligns with the conventions and philosophy of the Go programming language.

## Go Proverbs: The Philosophy of Go

Before diving into the specifics of coding, it's crucial to understand the philosophy that underpins the Go language. The "Go Proverbs" offer a glimpse into this philosophy. Here are some of the most important ones to keep in mind:

- **Don't communicate by sharing memory, share memory by communicating.** This is the core principle of Go's concurrency model. It favors passing data between goroutines via channels over using shared memory and locks.
- **Concurrency is not parallelism.** Concurrency is about dealing with lots of things at once. Parallelism is about doing lots of things at once. Go provides the tools for concurrency, and parallelism is a possible outcome.
- **Channels orchestrate; mutexes serialize.** Use channels to coordinate the work of goroutines. Use mutexes to protect access to shared data.
- **The bigger the interface, the weaker the abstraction.** Small, focused interfaces are more powerful and easier to use.
- **Make the zero value useful.** The zero value of a type should be a valid, usable state.
- **`interface{}` says nothing.** An empty interface provides no information about the type of data it holds. Use it sparingly.
- **`gofmt`'s style is no one's favorite, yet `gofmt` is everyone's favorite.** The Go community has embraced a single, automated formatting style. Use `gofmt` to format your code.
- **A little copying is better than a little dependency.** It's often better to copy a small amount of code than to introduce a new dependency.
- **Clear is better than clever.** Write code that is easy to understand.
- **Errors are values.** Treat errors as regular values that can be returned, passed around, and handled.

## Project Layout and Package-Oriented Design

A well-structured project is easier to understand and maintain. The following guidelines are based on the official Go documentation and Ardan Labs' package-oriented design principles.

### Standard Project Layout

- **`cmd/`**: This directory contains the `main` packages for your applications. Each subdirectory within `cmd/` should correspond to a single executable.
- **`internal/`**: This directory contains code that is specific to your application and should not be imported by other projects. The Go compiler enforces this. Business logic is often placed here, organized into subfolders by domain.
- **`pkg/`**: This directory is for reusable packages that are safe to be used by other projects.
- **`/` (root)**: For smaller projects, it's acceptable to have `main.go` in the root of the project.

### Package-Oriented Design Principles

- **Responsibility Isolation**: Each package should have a single, clearly defined purpose.
- **Dependency Minimization**: Packages should only depend on other packages when absolutely necessary. Avoid circular dependencies.
- **Dependency Inversion**: Use interfaces to abstract dependencies, especially for external packages.
- **Encapsulation**: Use Go's export rules (capitalized identifiers) to control which parts of a package are accessible from the outside, creating a clear API.
- **Packages Provide, Not Contain**: Focus on what a package *provides* in terms of functionality, rather than simply what it *contains*.

## Naming Conventions

- **Package Names**: Package names should be short, concise, and all lowercase. They should be single words.
- **Variable and Function Names**: Use `MixedCaps` or `mixedCaps` (CamelCase) for multi-word names.
- **Getters**: If a struct has a field `owner` (unexported), the getter method should be `Owner()` (exported), not `GetOwner()`.
- **Interface Names**: Interfaces that represent a single method are often named by the method name plus the "-er" suffix (e.g., `Reader`, `Writer`).

## Formatting and Style

- **`gofmt`**: Always use `gofmt` to format your code. This is the standard in the Go community and ensures consistency.
- **Line Length**: Go does not have a strict line length limit, but try to keep lines to a reasonable length to improve readability.
- **Mixed Caps**: Use `MixedCaps` for names, not underscores.

## Comments

- **Comment Sentences**: Comments should be complete sentences.
- **Exported Symbols**: Every exported symbol should have a doc comment.
- **Package Comments**: Each package should have a package comment, a block comment preceding the package clause.

## Functions and Methods

- **Multiple Return Values**: Go functions can return multiple values. This is often used to return a result and an error.
- **Named Result Parameters**: Avoid named result parameters in most cases. They can be confusing and are often unnecessary.
- **Naked Returns**: Avoid naked returns. Always be explicit about what you are returning.

## Interfaces

- **Small Interfaces**: Prefer small, focused interfaces. The bigger the interface, the weaker the abstraction.
- **`interface{}`**: Use the empty interface (`interface{}`) sparingly. It provides no type information.

## Error Handling

- **Errors are Values**: Treat errors as values. Return them from functions and handle them.
- **Explicit Error Handling**: Handle all errors explicitly. Do not discard them using the blank identifier (`_`).
- **Error Wrapping**: Use `fmt.Errorf` with the `%w` verb to wrap errors and provide context.
- **Don't Panic**: Do not use `panic` for normal error handling. `panic` is for exceptional, unrecoverable errors.

## Concurrency

- **Share Memory by Communicating**: Use channels to pass data between goroutines.
- **`select`**: Use the `select` statement to wait on multiple channel operations.
- **`context.Context`**: Use the `context` package to manage cancellation, deadlines, and other request-scoped values. Pass a `Context` as the first argument to functions that may be long-running or need to be cancellable.

## HTTP Services

- **`http.Handler`**: Use the `http.Handler` interface to implement HTTP handlers.
- **Dependencies as Arguments**: Pass dependencies to your handlers as arguments, rather than using global variables.
- **`NewServer` Constructor**: Use a `NewServer` constructor to set up your server and its dependencies.
- **Graceful Shutdown**: Implement graceful shutdown to allow your server to finish handling existing requests before shutting down.
- **Testing**: Use end-to-end testing for your HTTP services.

## Testing

- **Test Files**: Test files are named with the `_test.go` suffix.
- **Test Functions**: Test functions are named `TestXxx` and take a `*testing.T` as a parameter.
- **Table-Driven Tests**: Use table-driven tests to test multiple scenarios with the same test logic.
- **Testable Examples**: Provide testable examples to document and verify your code's behavior.

## Conclusion

This guide provides a starting point for writing idiomatic Go code. By following these principles, you will be able to generate code that is not only correct but also clean, maintainable, and in line with the best practices of the Go community. Remember that these are guidelines, not strict rules. The most important thing is to write code that is clear, simple, and easy to understand.
