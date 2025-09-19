# Backlog <!-- omit in toc -->

> Manage your project's backlog with Markdown in Git.
> Designed for seamless collaboration between developers and AI agents.

Backlog is a zero-configuration, offline-first task manager written in Go. Every task lives as a Markdown file inside your Git repository, and the Backlog CLI plus its MCP (Model Context Protocol) server keep those files consistent, searchable, and automation-friendly. Because the entire state travels with the repo, Backlog works anywhere your Git checkout does—laptops, containers, CI, or remote development environments.

<!--toc:start-->
- [Introduction](#introduction)
- [Why Backlog](#why-backlog)
- [Key Concepts](#key-concepts)
- [Quick Start](#quick-start)
  - [Install Backlog](#install-backlog)
  - [Create Your First Task](#create-your-first-task)
  - [Task Directory Resolution](#task-directory-resolution)
- [Configuration](#configuration)
- [AI Agent Integration](#ai-agent-integration)
  - [Available Tools](#available-tools)
  - [Registering AI Agents](#registering-ai-agents)
- [CLI Usage Examples](#cli-usage-examples)
  - [Everyday Task Management](#everyday-task-management)
  - [Hierarchical Task Structure](#hierarchical-task-structure)
  - [Advanced Task Creation](#advanced-task-creation)
- [Task File Anatomy](#task-file-anatomy)
- [Project Layout](#project-layout)
- [Development](#development)
- [Inspiration](#inspiration)
<!--toc:end-->

## Introduction

**THIS SECTION IS FOR HUMANS**

Large refactors and long-running efforts stretch an AI assistant's context window quickly. A reliable approach is to ask the agent to "make a plan for X, write it to a markdown file, and keep it updated as work progresses." Doing that by hand leaves the repository littered with `plan.md` and other ad-hoc files. Backlog fixes that by providing a first-class tasks folder plus an MCP server the agent can understand and trust. The result is a shared backlog that stays in sync whether tasks are updated by people or by automation.

You can browse the generated tasks for this repository under the [.backlog](./.backlog) folder. They were produced from the following prompt:

```
If you were to recreate this project from scratch, make a plan and break it down into tasks using backlog break down prompt.
Write that plan to a markdown file called "./docs/plan.md".
Check that the plan in ./docs/plan.md is consistent with the list of tasks in the backlog.
Add implementation plan to relevant tasks
All tasks should have at least one acceptance criteria.

Read the full instructions for backlog: ./internal/mcp/prompt.md
```

Backlog is registered as an MCP server with the `backlog mcp` command. Configuration examples live in [.gemini](./.gemini), [.claude](./.claude), and [.vscode](./.vscode).

While most of the codebase is hand-written, documentation, examples, and parts of the tests were produced with AI. In practice, AI agents deliver better output when the project already reflects the exact structure, libraries, and style you want. Once the scaffolding is there, collaboration becomes much easier.

Tools used while building this project (all paid tiers):

- gemini-2.5-pro
- claude-sonnet4
- github-copilot
- amp

> To instruct your AI tools to use `backlog`, read [./internal/mcp/prompt.md](./internal/mcp/prompt.md). It is also accessible via `backlog instructions > backlog_instructions.md`.

## Why Backlog

- **Markdown-first** – tasks are plain files that work with your existing Git workflow.
- **Offline & portable** – no database, no external service; everything travels with the repo.
- **Zero configuration** – run the binary and start creating tasks immediately.
- **AI-native** – the MCP server exposes a rich toolset so agents can plan, execute, and document work safely.
- **Git-aware** – automatic commits keep changes traceable when auto-commit is enabled.
- **Hierarchical** – parent/child IDs make breaking down work painless (e.g., `T01 → T01.01 → T01.01.01`).

## Key Concepts

- **Task ID**: Immutable, hierarchical identifier like `T03.01`. Generated automatically when you create tasks.
- **Task Files**: Stored in `.backlog/` as `T{ID}-{slugified-title}.md`.
- **TaskStore**: Backed by the filesystem via `afero`, making it easy to mock in tests.
- **Acceptance Criteria**: Represented as ordered checkboxes; managed via CLI or MCP tools.
- **MCP Tools**: Provide structured read/write access for agents (`task_create`, `task_edit`, etc.).

## Quick Start

### Install Backlog

```bash
# Build from source
git clone https://github.com/veggiemonk/backlog
cd backlog
go build .

# Or install directly
go install github.com/veggiemonk/backlog@latest

# Or pull the container
docker pull ghcr.io/veggiemonk/backlog:latest
```

You can also grab pre-built binaries from the [releases page](https://github.com/veggiemonk/backlog/releases).

### Create Your First Task

```bash
# Inside any Git repository
backlog create "Implement password reset" \
  -d "Users should be able to request a password reset link." \
  --ac "Email includes secure token" \
  --ac "Token expires after 30 minutes"

# See what you just created
backlog list
backlog view T01
```

You do not need an init step—Backlog creates the `.backlog/` directory on demand.

### Task Directory Resolution

Backlog stores tasks in a configurable "tasks folder" (default `.backlog`). Override it when needed.

- **Set the folder**
  - CLI flag: `--folder <path>` (relative or absolute)
  - Environment variable: `BACKLOG_FOLDER`
  - Default: `.backlog`

- **Resolution rules (applied to the chosen value)**
  - Absolute paths are used as-is.
  - Relative paths resolve with this precedence:
    1. If `<CWD>/<path>` exists, use it.
    2. Search parent directories for `<ancestor>/<path>`.
    3. If a Git repo exists, prefer `<gitRoot>/<path>`.
    4. Otherwise, fall back to `<CWD>/<path>` (created as needed).

- **Container tips**
  - If `.git` is absent in the container, upward search still works.
  - For deterministic setups, mount the tasks directory and set `BACKLOG_FOLDER` to the absolute mount point (or pass `--folder`).

## Configuration

Backlog reads configuration from CLI flags (highest precedence) and environment variables.

| Setting | Flag | Env Var | Default | Description |
| ------- | ---- | ------- | ------- | ----------- |
| Tasks Directory | `--folder` | `BACKLOG_FOLDER` | `.backlog` | Location for task files |
| Auto Commit | `--auto-commit` | `BACKLOG_AUTO_COMMIT` | `true` | Automatically commit task changes |
| Log Level | `--log-level` | `BACKLOG_LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| Log Format | `--log-format` | `BACKLOG_LOG_FORMAT` | `text` | `text` or `json` |
| Log File | `--log-file` | `BACKLOG_LOG_FILE` | *(stderr)* | Destination for logs |

Examples:

```bash
# Flags
backlog --folder /backlog --log-level debug --auto-commit false list

# Environment variables
export BACKLOG_FOLDER=/backlog
export BACKLOG_LOG_LEVEL=debug
export BACKLOG_AUTO_COMMIT=false
backlog list

# Flags override environment variables
export BACKLOG_FOLDER=/env/path
backlog --folder /flag/path list  # uses /flag/path

# Logging
backlog --log-format json --log-level debug list
BACKLOG_LOG_FILE=/tmp/backlog.log backlog list
```

Notes:

- All environment variables use the `BACKLOG_` prefix.
- Boolean env values should be `true`/`false` strings.
- When `--log-file` is omitted, logs go to stderr.

## AI Agent Integration

Backlog ships with an MCP server that mirrors CLI functionality for AI agents. Agents can plan, create, and update tasks without ever touching the raw Markdown files.

### Available Tools

- `task_create`: Create new tasks with full metadata.
- `task_list`: Filter and page through tasks.
- `task_view`: Retrieve a single task.
- `task_edit`: Update titles, descriptions, status, acceptance criteria, and more.
- `task_search`: Query tasks by content.
- `task_archive`: Archive tasks without deleting the underlying files.

### Registering AI Agents

```bash
# HTTP transport (default port 8106)
backlog mcp --http
backlog mcp --http --port 9000

# STDIO transport
backlog mcp
```

See `.gemini/`, `.claude/`, and `.vscode/` for real configuration snippets. For full agent-facing guidance, reference [internal/mcp/prompt.md](./internal/mcp/prompt.md).

## CLI Usage Examples

### Everyday Task Management

```bash
# List everything
backlog list

# Filter by status
backlog list --status todo
backlog list --status in-progress
backlog list --status done

# View or edit a single task
backlog view T01.02
backlog edit T01 --status in-progress --assignee "alex"
```

### Hierarchical Task Structure

```bash
# Parent task
backlog create "Implement User Authentication"
# → T01-implement_user_authentication.md

# Subtask
backlog create "Add Google OAuth login" -p T01
# → T01.01-add_google_oauth_login.md

# Sub-subtask
backlog create "OAuth token validation" -p T01.01
# → T01.01.01-oauth_token_validation.md
```

### Advanced Task Creation

```bash
backlog create "Build reporting feature" \
  -d "Create monthly performance reports in PDF format." \
  -a drew \
  -l "feature,frontend,backend" \
  --priority high \
  --ac "Report generation logic is accurate" \
  --ac "Users can select date range" \
  --ac "PDF export works correctly" \
  -p 23
```

## Task File Anatomy

Tasks are Markdown files with YAML frontmatter stored under `.backlog/`.

```markdown
---
id: "01.02.03"
title: "Implement OAuth integration"
status: "todo"
parent: "01.02"
assigned: ["alex", "jordan"]
labels: ["feature", "auth", "backend"]
priority: "high"
created_at: 2024-01-01T00:00:00Z
updated_at: 2024-01-01T00:00:00Z
---

## Description
Integrate OAuth authentication with Google and GitHub providers.

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Google OAuth integration works
- [ ] #2 GitHub OAuth integration works
- [x] #3 OAuth scope validation implemented
<!-- AC:END -->

## Implementation Plan
1. Setup OAuth app registrations
2. Implement OAuth flow handlers
3. Add token validation
4. Write integration tests

## Implementation Notes
- Use oauth2 package for Go
- Store tokens securely
- Handle refresh token rotation
```

File naming convention examples:

- `T01-implement_user_auth.md`
- `T01.01-setup_oauth.md`
- `T01.01.01-google_oauth.md`

## Project Layout

```
.backlog/                     # Task storage directory
├── T01-user_auth.md         # Root task (ID: T01)
├── T01.01-oauth_setup.md    # Subtask (ID: T01.01)
├── T01.02-password_reset.md # Sibling subtask (ID: T01.02)
└── T02-frontend_redesign.md # Another root task (ID: T02)
```

## Development

```bash
make build   # Compile the CLI
make test    # Run the full test suite
make docs    # Generate Cobra command docs
make lint    # Run linters
```

## Inspiration

Backlog draws inspiration from:

- [Backlog.md](https://github.com/MrLesk/Backlog.md)
- [TaskWing](https://github.com/josephgoksu/TaskWing)

For MCP-specific inspiration, see [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) and the production-grade examples in [GoogleCloudPlatform/gke-mcp](https://github.com/GoogleCloudPlatform/gke-mcp). Backlog aims to stay simpler and laser-focused on AI-friendly workflows.
