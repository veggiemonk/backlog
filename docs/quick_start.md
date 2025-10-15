# Quick Start



### Installation

```bash
# Build from source
git clone https://github.com/veggiemonk/backlog
cd backlog
go build .

# Or install directly
go install github.com/veggiemonk/backlog@latest

# Or in a container
docker pull ghcr.io/veggiemonk/backlog:latest
```

You can also download the binary from the [release page](https://github.com/veggiemonk/backlog/releases).

#### Using the Container

The backlog container is designed to work with your local project directory. Here are common usage patterns:

```bash
# Basic usage - mount current directory and set backlog folder
docker run --rm -it -v $(pwd):/data -e BACKLOG_FOLDER=/data/.backlog ghcr.io/veggiemonk/backlog list

# Create a task
docker run --rm -it -v $(pwd):/data -e BACKLOG_FOLDER=/data/.backlog ghcr.io/veggiemonk/backlog create "Fix bug in authentication"

# View a specific task
docker run --rm -it -v $(pwd):/data -e BACKLOG_FOLDER=/data/.backlog ghcr.io/veggiemonk/backlog view T01

# Edit a task
docker run --rm -it -v $(pwd):/data -e BACKLOG_FOLDER=/data/.backlog ghcr.io/veggiemonk/backlog edit T01 --status "in-progress"

# Start MCP server (for AI integration)
docker run --rm -it -v $(pwd):/data -e BACKLOG_FOLDER=/data/.backlog -p 8106:8106 ghcr.io/veggiemonk/backlog mcp --http --port 8106
```

**Container Tips:**

- Mount your project directory to `/data` for file persistence
- Set `BACKLOG_FOLDER=/data/.backlog` to store tasks in your project
- Use `-p 8106:8106` when running the MCP server to expose the port
- The `--rm` flag removes the container after execution
- All CLI commands work the same way in the container

### Initialize Your Project

No initialization is needed.

### Task Directory Resolution

Backlog stores tasks in a directory referred to as the "tasks folder". By default this is `.backlog`, but you can override it.

#### How to set the folder

- CLI flag: `--folder <path>` (relative or absolute)
- Environment variable: `BACKLOG_FOLDER` (used when set)
- Default: `.backlog`

#### Resolution rules (applied to the chosen value)

- Absolute path: used as-is.
- Relative path: resolved with this precedence:
  - If `<CWD>/<path>` exists, use it.
  - Search parent directories; if `<ancestor>/<path>` exists, use it.
  - If a git repository is detected, use `<gitRoot>/<path>`.
  - Otherwise, fall back to `<CWD>/<path>` (created on demand).

#### Container tips

- If your container does not include the `.git` directory, the resolver still works using the upward search and CWD fallback.
- For predictable behavior, mount your tasks directory and set `BACKLOG_FOLDER` to its absolute mount point, or pass `--folder` with an absolute path.

##
