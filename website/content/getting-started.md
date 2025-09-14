---
title: "Getting Started"
description: "Installation and basic usage guide for Backlog"
weight: 20
---

# Getting Started

This guide will help you get up and running with Backlog quickly.

## Installation

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

## First Steps

{{< alert >}}
No initialization is required! Backlog works immediately in any directory.
{{< /alert >}}

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

{{< details "Simple Task Creation" >}}
```bash
# Task with description and priority
backlog create "Implement user authentication" \
  -d "Add login and registration functionality" \
  --priority "high"
```
{{< /details >}}

{{< details "Task with Assignees and Labels" >}}
```bash
backlog create "Update dependencies" \
  -a "alex" -a "jordan" \
  -l "maintenance,backend,security"
```
{{< /details >}}

{{< details "Task with Acceptance Criteria" >}}
```bash
backlog create "Build reporting feature" \
  -d "Create monthly performance reports in PDF format" \
  --ac "Report generation logic is accurate" \
  --ac "Users can select date range" \
  --ac "PDF export works correctly"
```
{{< /details >}}

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

{{< tabpane >}}

#### Update Status
```bash
backlog edit T01 --status "in-progress"
```

#### Add Assignee
```bash
backlog edit T01 --add-assignee "alex"
```

<--->

#### Add Labels
```bash
backlog edit T01 --add-label "urgent" --add-label "frontend"
```

#### Complete Acceptance Criteria
```bash
backlog edit T01 --check-ac 1,3
```

{{< /tabpane >}}

### Filtering and Searching

{{< tabpane >}}
{{< tab "By Status" >}}
```bash
backlog list --status "todo"
backlog list --status "in-progress"
```
{{< /tab >}}
{{< tab "By Assignee" >}}
```bash
backlog list --assignee "alex"
```
{{< /tab >}}
{{< tab "By Parent" >}}
```bash
# Show subtasks
backlog list --parent "T01"
```
{{< /tab >}}
{{< tab "Search Content" >}}
```bash
backlog search "authentication"
```
{{< /tab >}}
{{< /tabpane >}}

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

{{< details "Task File Example" >}}
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
{{< /details >}}

## Git Integration

{{< alert >}}
Backlog automatically commits changes to your Git repository:
{{< /alert >}}

```bash
# After creating a task
git log -1 --oneline
# create: T01 Set up project documentation

# After updating a task
git log -1 --oneline
# update: T01 moved to in-progress, assigned alex
```

## Next Steps

[**CLI Reference**]({{< relref "cli" >}}){: .btn .btn-primary}
[**AI Integration**]({{< relref "ai-integration" >}}){: .btn .btn-primary}
[**Development Guide**]({{< relref "development" >}}){: .btn .btn-primary}