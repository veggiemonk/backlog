---
title: "Backlog"
description: "Zero-configuration task manager for developers and AI agents"
geekdocNav: false
geekdocAlign: center
geekdocAnchor: false
---

[**Get Started**](/getting-started){: .btn .btn-primary .btn-lg}
[**GitHub**](https://github.com/veggiemonk/backlog){: .btn .btn-outline-primary .btn-lg}

## Zero-configuration task manager for developers and AI agents

Backlog is a task manager written in Go where tasks live inside a Git repository. It leverages plain Markdown files for task storage and a comprehensive command-line interface (CLI) for interaction. This design makes it exceptionally well-suited for AI agents thanks to its MCP (Model Context Protocol) integration.

<div class="row">
  <div class="col-md-6">
    <h2>Features</h2>
    <ul>
      <li><strong>Task Management</strong>: Create, edit, list, and view tasks with rich metadata</li>
      <li><strong>Hierarchical Structure</strong>: Support for parent-child-grandchild task relationships (T01 → T01.01 → T01.01.01)</li>
      <li><strong>Search & Filter</strong>: Find tasks by content, status, parent relationships, and labels</li>
      <li><strong>AI-Friendly</strong>: MCP server integration for seamless AI agent collaboration</li>
    </ul>
  </div>
  <div class="col-md-6">
    <h2>Benefits</h2>
    <ul>
      <li><strong>Git Integration</strong>: Tasks are stored as Markdown files with automatic Git commits</li>
      <li><strong>Offline-First</strong>: Works completely offline with local Git repository storage</li>
      <li><strong>Portable</strong>: Entire project state contained within the Git repository</li>
      <li><strong>Zero Configuration</strong>: No setup files or databases required</li>
    </ul>
  </div>
</div>

## Quick Start

### Installation

{{< tabpane >}}
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
{{< /tabpane >}}

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

{{< alert title="Note" >}}
**For Humans**: The context window deteriorates rapidly on large-scale projects. A workaround is to ask AI agents to "make a plan for X, write it to a markdown file and keep it updated with the ongoing tasks". This technique has worked incredibly well, making refactoring and other significant code changes more resilient to failures.
{{< /alert >}}

{{< alert title="Tip" color="info" >}}
**For AI Agents**: Backlog provides a structured MCP server that AI tools can understand and trust to handle task management details, providing a much better experience when using AI tools for complex tasks.
{{< /alert >}}