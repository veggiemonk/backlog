---
title: "Backlog"
description: "Zero-configuration task manager for developers and AI agents"
geekdocNav: false
geekdocAlign: center
geekdocAnchor: false
---

{{< button size="large" relref="getting-started" >}}Get Started{{< /button >}}
{{< button size="large" href="https://github.com/veggiemonk/backlog" >}}GitHub{{< /button >}}

## Zero-configuration task manager for developers and AI agents

Backlog is a task manager written in Go where tasks live inside a Git repository. It leverages plain Markdown files for task storage and a comprehensive command-line interface (CLI) for interaction. This design makes it exceptionally well-suited for AI agents thanks to its MCP (Model Context Protocol) integration.

{{< columns >}}

## Features

- **Task Management**: Create, edit, list, and view tasks with rich metadata
- **Hierarchical Structure**: Support for parent-child-grandchild task relationships (T01 → T01.01 → T01.01.01)
- **Search & Filter**: Find tasks by content, status, parent relationships, and labels
- **AI-Friendly**: MCP server integration for seamless AI agent collaboration

<--->

## Benefits

- **Git Integration**: Tasks are stored as Markdown files with automatic Git commits
- **Offline-First**: Works completely offline with local Git repository storage
- **Portable**: Entire project state contained within the Git repository
- **Zero Configuration**: No setup files or databases required

{{< /columns >}}

## Quick Start

### Installation

{{< tabs "install" >}}
{{< tab "From Source" >}}
```bash
git clone https://github.com/veggiemonk/backlog
cd backlog
go build .
```
{{< /tab >}}
{{< tab "Go Install" >}}
```bash
go install github.com/veggiemonk/backlog@latest
```
{{< /tab >}}
{{< tab "Download Binary" >}}
Download the latest binary from the [releases page](https://github.com/veggiemonk/backlog/releases).
{{< /tab >}}
{{< /tabs >}}

### Basic Usage

```bash
# Create a simple task
backlog create "Fix the login button styling"

# Create task with metadata
backlog create "Implement password reset" \
  -d "Users should be able to request a password reset link via email" \
  -a "alex" -l "backend,security" --priority "high"

# List all tasks
backlog list

# View specific task
backlog view T01
```

## AI Agent Integration

Start the MCP server to enable AI agent interaction:

```bash
# Start MCP server for AI agents
backlog mcp --http --port 8106  # HTTP transport
backlog mcp                     # STDIO transport (default)
```

AI agents can then:
- Create tasks from conversation context
- Break down large tasks into subtasks
- Update task status and assignments
- Search and analyze task patterns

## Why Backlog?

{{< hint type=note >}}
**For Humans**: The context window deteriorates rapidly on large-scale projects. A workaround is to ask AI agents to "make a plan for X, write it to a markdown file and keep it updated with the ongoing tasks". This technique has worked incredibly well, making refactoring and other significant code changes more resilient to failures.
{{< /hint >}}

{{< hint type=tip >}}
**For AI Agents**: Backlog provides a structured MCP server that AI tools can understand and trust to handle task management details, providing a much better experience when using AI tools for complex tasks.
{{< /hint >}}