---
layout: page
title: Getting Started
nav_order: 3
---

# Getting Started

This guide will help you get up and running with Backlog quickly.

## First Steps

No initialization is required! Backlog works immediately in any directory.

### Create Your First Task

```bash
backlog create "Set up project documentation"
```

This creates a task file at `.backlog/T01-set_up_project_documentation.md`.

### View Your Tasks

```bash
backlog list
```

### View Task Details

```bash
backlog view T01
```

## Working with Tasks

### Creating Tasks with Metadata

```bash
# Task with description and priority
backlog create "Implement user authentication" \
  -d "Add login and registration functionality" \
  --priority "high"

# Task with assignees and labels
backlog create "Update dependencies" \
  -a "alex" -a "jordan" \
  -l "maintenance,backend,security"

# Task with acceptance criteria
backlog create "Build reporting feature" \
  -d "Create monthly performance reports in PDF format" \
  --ac "Report generation logic is accurate" \
  --ac "Users can select date range" \
  --ac "PDF export works correctly"
```

### Hierarchical Tasks

Create subtasks by specifying a parent:

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

### Managing Tasks

```bash
# Update task status
backlog edit T01 --status "in-progress"

# Add assignee
backlog edit T01 --add-assignee "alex"

# Add labels
backlog edit T01 --add-label "urgent" --add-label "frontend"

# Mark acceptance criteria as complete
backlog edit T01 --check-ac 1,3
```

### Filtering and Searching

```bash
# Filter by status
backlog list --status "todo"
backlog list --status "in-progress"

# Filter by assignee
backlog list --assignee "alex"

# Filter by parent (show subtasks)
backlog list --parent "T01"

# Search content
backlog search "authentication"
```

## File Structure

Tasks are stored as Markdown files in the `.backlog/` directory:

```
.backlog/
├── T01-implement_user_auth.md         # Root task
├── T01.01-setup_oauth.md              # Subtask
├── T01.01.01-google_oauth.md          # Sub-subtask
├── T02-frontend_redesign.md           # Another root task
└── archived/                          # Archived tasks
    └── T03-completed_feature.md
```

Each task file contains YAML frontmatter with metadata and Markdown content:

```markdown
---
id: "01.02"
title: "Setup OAuth integration"
status: "todo"
parent: "01"
assigned: ["alex"]
labels: ["feature", "auth"]
priority: "high"
created_at: 2024-01-01T00:00:00Z
updated_at: 2024-01-01T00:00:00Z
---

## Description

Set up OAuth integration with Google and GitHub providers.

## Acceptance Criteria

<!-- AC:BEGIN -->
- [ ] #1 Google OAuth integration works
- [ ] #2 GitHub OAuth integration works
- [x] #3 OAuth scope validation implemented
<!-- AC:END -->

## Implementation Notes

- Use oauth2 package for Go
- Store tokens securely
```

## Git Integration

Backlog automatically commits changes to your Git repository:

```bash
# After creating a task
git log -1 --oneline
# create: T01 Set up project documentation

# After updating a task
git log -1 --oneline
# update: T01 moved to in-progress, assigned alex
```

## Next Steps

- Learn about [AI Integration](ai-integration.md) to use Backlog with AI agents
- Explore the complete [CLI Reference](cli/backlog.md) for all available commands
- Check out the [Development Guide](development/index.md) if you want to contribute
