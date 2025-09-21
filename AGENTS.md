# AI Agent Instructions for the `backlog` Repository

This document provides a comprehensive and unified set of instructions for AI agents working on this project.

## 1. Project Overview

- **Project Name**: `backlog`
- **Description**: A Go-based, offline-first, portable task management CLI that stores tasks as Markdown files within a Git repository.
- **Key Feature**: Designed for seamless integration with AI agents via CLI commands.

---

## 2. Core Architecture & Concepts

To work effectively, you must understand these core components:

- **`Task`**: The fundamental data unit. A task is a Markdown file with YAML frontmatter.
- **`TaskStore`**: The storage abstraction layer. The primary implementation is `FileTaskStore`, which uses the `afero` library for filesystem operations (allowing for in-memory mocking in tests).
- **`TaskID`**: A hierarchical, dot-notation identifier (e.g., `T01`, `T01.01`). IDs are auto-generated and should be treated as immutable.
- **Storage Location**: All task files are stored in the `.backlog/` directory.
- **File Naming**: Task files follow the convention: `T{ID}-{slugified-title}.md`.
- **Directory Structure**:
    - `internal/core/`: Contains all core business logic (creating, listing, updating tasks).
    - `internal/cmd/`: Defines the CLI commands using the Cobra framework.
    - `internal/commit/`: Handles automatic Git integration.
    - `internal/mcp/`: Implements the MCP server for AI agent tool-based interaction.
    - `internal/logging/`: Logging for this project.

---

## 3. Development Workflow

Follow these steps for reading, analyzing, and editing code.

### 3.1. Standard Build & Test Commands

Use `make` for simplicity.

- **Build**: `make build`
- **Run all tests**: `make test`
- **Run tests for a specific package**: `go test ./internal/core`
- **Lint**: `make lint`
- **Generate Docs**: `make docs`
- **Install**: `make install`

### 3.2. Go-Specific Analysis & Editing Workflow (`gopls`)

You MUST use the `gopls` tools for code intelligence.

**Reading & Understanding Code:**

1.  **Workspace Overview**: Start with `go_workspace` to understand the project structure.
2.  **Find Symbols**: Use `go_search` with a fuzzy query to locate types, functions, or variables.
3.  **File Context**: After reading a Go file, immediately use `go_file_context` to understand its dependencies within the same package.
4.  **Package API**: Use `go_package_api` to understand the public API of any package (internal or third-party).

**Editing Code:**

1.  **Understand First**: Follow the reading workflow above before any modification.
2.  **Find References**: Before changing a symbol, you MUST use `go_symbol_references` to find all its usages and assess the impact.
3.  **Make Edits**: Perform all necessary code changes.
4.  **Check for Errors**: After every edit, you MUST run `go_diagnostics` on the changed files.
5.  **Fix Errors**: If `go_diagnostics` reports errors, fix them. Review and apply suggested quick fixes if they are correct.
6.  **Run Tests**: Once `go_diagnostics` is clean, run the relevant tests using `go test [packagePath...]`.

---

## 4. Critical Rules & Safety Guidelines

**Non-negotiable rules for interacting with the task system:**

1.  **NEVER EDIT TASK FILES DIRECTLY**: Do not write to, modify, or delete any file in the `.backlog/` directory. All task manipulations MUST go through the `backlog` CLI commands. Direct edits will corrupt metadata and break the system.
2.  **DO NOT DELETE TASK FILES**: Never delete task markdown files. Use the `backlog archive` command instead.
3.  **ALWAYS USE THE `backlog` TOOL**: For any operation related to task management (create, list, view, edit, search, archive), you MUST use the `backlog` CLI tool.

---

## 5. Planning and Task Management (`backlog`)

FULLY READ THE INSTRUCTIONS FOR BACKLOG CLI [prompt.md](./internal/mcp/prompt.md)

---

## 6. Code Style & Architectural Patterns

- **Error Handling**: Wrap errors with context using `fmt.Errorf("context: %w", err)`.
- **Dependency Injection**: Business logic is implemented as methods on the `FileTaskStore`. The `FileTaskStore` is created with an `afero.Fs` instance, allowing for dependency injection.
- **TaskID Parsing**: Always use `core.ParseTaskID()` to handle user-provided task IDs, as it supports flexible formats (e.g., "T1.2", "1.2").
- **Testing**: Use `afero.NewMemMapFs()` to create an in-memory filesystem for tests.
- **Git Integration**: Task operations trigger automatic Git commits. Ensure the repository is a valid Git repository.

---

## 7. AI Agent 

### 7.1. Specialist Agents

- **`go-task-manager-reviewer`**: An expert agent for reviewing Go code related to this project.
- **When to use**: Use this agent when you need an expert review of new features, architectural changes, or refactoring, especially concerning task management logic, storage, or CLI commands.
