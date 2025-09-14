---
title: "Development"
description: "Architecture, contribution guidelines, and technical documentation"
weight: 50
geekdocCollapseSection: true
---

# Development

Resources for developers who want to understand, contribute to, or extend the Backlog project.

## Overview

This section contains technical documentation about the Backlog project's architecture, implementation details, and development processes.

{{< columns >}}

### Technical Documentation
- [**Technical Analysis**](analysis) - Comprehensive analysis of the codebase, architecture, and implementation patterns
- [**Project Plan**](plan) - Detailed recreation plan with development phases and acceptance criteria
- [**Architecture Guide**](architecture) - System design and architectural decisions

<--->

### Contributing Resources
- [**Getting Started**]({{< relref "getting-started" >}}) - Basic usage and installation
- [**GitHub Repository**](https://github.com/veggiemonk/backlog) - Source code and issues
- [**CLI Reference**]({{< relref "cli" >}}) - Complete command documentation

{{< /columns >}}

## Contributing

{{< hint type=tip >}}
Backlog is designed to be simple and focused. Before contributing:
{{< /hint >}}

1. **Read the Technical Analysis** to understand the current architecture
2. **Review the Project Plan** to see the intended direction
3. **Check existing issues** on [GitHub](https://github.com/veggiemonk/backlog/issues)

### Development Environment

{{< tabs "dev-setup" >}}
{{< tab "Quick Start" >}}
```bash
# Clone the repository
git clone https://github.com/veggiemonk/backlog
cd backlog

# Build the project
make build
```
{{< /tab >}}
{{< tab "Testing" >}}
```bash
# Run tests
make test

# Run with coverage
make cover
```
{{< /tab >}}
{{< tab "Documentation" >}}
```bash
# Generate CLI documentation
make docs

# Serve documentation locally
hugo server
```
{{< /tab >}}
{{< tab "Quality" >}}
```bash
# Lint code
make lint

# Format code
go fmt ./...
```
{{< /tab >}}
{{< /tabs >}}

### Architecture Principles

{{< columns >}}

#### Design Philosophy
- **File-based storage**: Tasks as Markdown files with YAML frontmatter
- **Interface-driven design**: TaskStore interface with filesystem abstraction
- **Git integration**: Automatic commits for all operations
- **Zero configuration**: Works immediately without setup

<--->

#### Technical Approach
- **MCP integration**: AI agent compatibility through Model Context Protocol
- **Clean separation**: Core business logic separate from CLI and MCP server
- **Dependency injection**: Testable design with afero filesystem mocking
- **Modern Go**: Go 1.21+ with contemporary idioms and patterns

{{< /columns >}}

### Code Quality Standards

{{< expand "Testing Requirements" >}}
- Comprehensive test coverage using afero filesystem mocking
- Unit tests for all core business logic
- Integration tests for CLI commands
- MCP server functionality tests
- Tests run against in-memory filesystems for isolation
{{< /expand >}}

{{< expand "Code Style" >}}
- Go 1.21+ with modern idioms
- Clean separation of concerns (core, CLI, MCP server)
- Interface-driven design for testability
- Consistent error handling patterns
- Meaningful variable and function names
{{< /expand >}}

{{< expand "Documentation" >}}
- All public APIs documented with Go doc comments
- CLI commands auto-generated with Cobra
- Architecture decisions recorded in technical analysis
- Examples provided for common use cases
{{< /expand >}}

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
├── content/               # Documentation (Hugo)
├── .backlog/              # Example backlog tasks
└── dist/                  # Release artifacts
```

## Design Philosophy

{{< hint type=important title="Core Principles" >}}
**Simplicity**: Handle fewer use cases well rather than many use cases poorly.

**AI-Friendly**: Designed from the ground up to work seamlessly with AI agents through MCP integration.

**Developer-Centric**: Built by developers, for developers, with a focus on transparency and control.

**Offline-First**: Everything works locally without external dependencies or cloud services.
{{< /hint >}}

## Getting Help

{{< columns >}}

### Community Resources
- **Issues**: [GitHub Issues](https://github.com/veggiemonk/backlog/issues)
- **Discussions**: [GitHub Discussions](https://github.com/veggiemonk/backlog/discussions)
- **Source Code**: [GitHub Repository](https://github.com/veggiemonk/backlog)

<--->

### Documentation
- **User Guide**: [Getting Started]({{< relref "getting-started" >}})
- **AI Integration**: [MCP Setup]({{< relref "ai-integration" >}})
- **API Reference**: [CLI Commands]({{< relref "cli" >}})

{{< /columns >}}