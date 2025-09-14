---
layout: page
title: Development
nav_order: 5
has_children: true
---

# Development

Resources for developers who want to understand, contribute to, or extend the Backlog project.

## Overview

This section contains technical documentation about the Backlog project's architecture, implementation details, and development processes.

## Contents

- [**Technical Analysis**](analysis.html) - Comprehensive analysis of the codebase, architecture, and implementation patterns
- [**Project Plan**](plan.html) - Detailed recreation plan with development phases and acceptance criteria
- [**Architecture**](architecture.html) - System design and architectural decisions

## Contributing

Backlog is designed to be simple and focused. Before contributing:

1. **Read the Technical Analysis** to understand the current architecture
2. **Review the Project Plan** to see the intended direction
3. **Check existing issues** on [GitHub](https://github.com/veggiemonk/backlog/issues)

### Development Environment

```bash
# Clone the repository
git clone https://github.com/veggiemonk/backlog
cd backlog

# Build the project
make build

# Run tests
make test

# Generate CLI documentation
make docs
```

### Architecture Principles

- **File-based storage**: Tasks as Markdown files with YAML frontmatter
- **Interface-driven design**: TaskStore interface with filesystem abstraction
- **Git integration**: Automatic commits for all operations
- **MCP integration**: AI agent compatibility through Model Context Protocol
- **Zero configuration**: Works immediately without setup

### Code Quality

- Go 1.21+ with modern idioms
- Comprehensive test coverage using afero filesystem mocking
- Clean separation of concerns (core, CLI, MCP server)
- Dependency injection for testability

### Testing

```bash
# Run all tests
make test

# Run with coverage
make cover

# Lint code
make lint
```

The project uses afero filesystem abstraction for testability, allowing tests to run against in-memory filesystems.

## Project Structure

```
backlog/
├── main.go                 # Entry point
├── internal/
│   ├── core/              # Business logic and data structures
│   ├── cmd/               # CLI commands (Cobra)
│   ├── mcp/               # MCP server implementation
│   ├── commit/            # Git integration
│   └── logging/           # Logging configuration
├── docs/                  # Documentation (this site)
├── .backlog/              # Example backlog tasks
└── dist/                  # Release artifacts
```

## Design Philosophy

**Simplicity**: Handle fewer use cases well rather than many use cases poorly.

**AI-Friendly**: Designed from the ground up to work seamlessly with AI agents through MCP integration.

**Developer-Centric**: Built by developers, for developers, with a focus on transparency and control.

**Offline-First**: Everything works locally without external dependencies or cloud services.

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/veggiemonk/backlog/issues)
- **Discussions**: [GitHub Discussions](https://github.com/veggiemonk/backlog/discussions)
- **Code**: [Source Code](https://github.com/veggiemonk/backlog)