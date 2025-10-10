# Usage Examples



### AI Integration (MCP Server)

```bash
# Start MCP server for AI agents (backlog mcp --http transport)
backlog mcp --http  # default port 8106
backlog mcp --http --port 8106 # specify the port

# Start MCP server (stdio transport)
backlog mcp
```

See this repository [.gemini](https://github.com/veggiemonk/backlog/tree/main/.gemini) or [.claude](https://github.com/veggiemonk/backlog/tree/main/.claude) for example configurations.

### Basic Task Creation

```bash
# Simple task
backlog create "Fix the login button styling"

# Task with description
backlog create "Implement password reset" \
  -d "Users should be able to request a password reset link via email"

# Task with assignees and labels
backlog create "Update dependencies" \
  -a "alex" -a "jordan" \
  -l "maintenance,backend,security" \
  --priority "high"
```

### Hierarchical Task Structure

```bash
# Create parent task
backlog create "Implement User Authentication"
# → Creates T01-implement_user_authentication.md

# Create subtask
backlog create "Add Google OAuth login" -p "T01"
# → Creates T01.01-add_google_oauth_login.md

# Create sub-subtask
backlog create "OAuth token validation" -p "T01.01"
# → Creates T01.01.01-oauth_token_validation.md
```

### Advanced Task Creation

```bash
# Comprehensive task with acceptance criteria
backlog create "Build reporting feature" \
  -d "Create monthly performance reports in PDF format" \
  -a "drew" \
  -l "feature,frontend,backend" \
  --priority "high" \
  --ac "Report generation logic is accurate" \
  --ac "Users can select date range" \
  --ac "PDF export works correctly" \
  -p "23"
```

### Task Dependencies

You can define dependencies between tasks using the `--deps` flag. A task cannot be started until its dependencies are complete.

```bash
# Create a task that depends on another
backlog create "Deploy to production" --deps "T10"

# Create a task with multiple dependencies
backlog create "Final release" --deps "T10,T11"

# Add a dependency to an existing task
backlog edit T12 --deps "T13"
```

### Task Management

```bash
# List all tasks
backlog list

# Filter by status
backlog list --status "todo"
backlog list --status "in-progress"
backlog list --status "done"

# Filter by parent (show subtasks)
backlog list --parent "T01"

# Pagination for large task lists
backlog list --limit 10                         # First 10 tasks
backlog list --limit 5 --offset 10              # Tasks 11-15
backlog list --status "todo" --limit 3          # First 3 todo tasks

# Search with pagination
backlog list --query "api" --limit 5                  # First 5 API-related tasks
backlog list --query "bug" --limit 3 --offset 5       # Search results 6-8

# View specific task
backlog view T01.02

# Edit task
backlog edit T01 --status "in-progress" --assigned "alex"
```

### Conflict Management

When working with Git branches, task ID conflicts can occur when multiple branches create tasks with the same IDs. Backlog provides automatic detection and resolution capabilities:

```bash
# Detect ID conflicts
backlog doctor                           # Text output
backlog doctor --json                    # JSON output

# Automatically fix conflicts
backlog doctor --fix                     # Fix using chronological strategy (default)
backlog doctor --fix --strategy=auto    # Fix using auto-renumber strategy
backlog doctor --fix --strategy=manual  # Create manual resolution plan
backlog doctor --fix --dry-run          # Preview changes without applying
```

**Conflict Types Detected:**

- **Duplicate IDs**: Same ID appears in multiple task files
- **Orphaned Children**: Tasks reference non-existent parent IDs
- **Invalid Hierarchy**: Parent-child relationships don't match ID structure

**Resolution Strategies:**

- **Chronological**: Keeps older tasks unchanged, renumbers newer conflicting tasks
- **Auto**: Automatically renumbers conflicting IDs using available ID space
- **Manual**: Creates a resolution plan requiring manual intervention

**Git Integration:**

Backlog automatically checks for conflicts during Git operations and can auto-resolve them during merges to maintain task integrity.

```bash
# Container usage for conflict management
docker run --rm -it -v $(pwd):/data -e BACKLOG_FOLDER=/data/.backlog ghcr.io/veggiemonk/backlog doctor
docker run --rm -it -v $(pwd):/data -e BACKLOG_FOLDER=/data/.backlog ghcr.io/veggiemonk/backlog doctor --fix --dry-run
```

##
