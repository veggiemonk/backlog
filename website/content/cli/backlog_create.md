---
title: "backlog create"
description: "Creates a new task in the backlog."
weight: 2
---

## backlog create

Create a new task.

### Synopsis

Creates a new task in the backlog. You can specify everything from a simple title to a detailed set of requirements, assignees, and dependencies.

```
backlog create <title> [flags]
```

{{< details "Show Examples" >}}
{{< tabpane >}}
{{< tab "Basic" >}}
##### 1. Basic Task Creation
The simplest way to create a task is with just a title.

```bash
backlog create "Fix the login button styling"
```

##### 2. Task with a Description
Use the `-d` or `--description` flag to add more detail.

```bash
backlog create "Implement password reset" -d "Users should be able to request a password reset link via their email."
```
{{< /tab >}}
{{< tab "Metadata" >}}
##### 3. Assigning a Task
Assign a task to one or more people using `-a` or `--assigned`.

```bash
# Assign to a single person
backlog create "Design the new dashboard" -a "alex"

# Assign to multiple people
backlog create "Code review for payment gateway" -a "jordan" -a "casey"
```

##### 4. Adding Labels
Use the `-l` or `--labels` flag to categorize the task.

```bash
backlog create "Update third-party dependencies" -l "bug,backend,security"
```

##### 5. Setting a Priority
Specify priority with `--priority`. Default is `medium`.

```bash
backlog create "Hotfix: Production database is down" --priority "high"
```
{{< /tab >}}
{{< tab "Advanced" >}}
##### 6. Defining Acceptance Criteria
Use `--ac` multiple times to define completion criteria.

```bash
backlog create "Develop user profile page" \
  --ac "User can view their profile info." \
  --ac "User can upload a new profile picture." \
  --ac "Page is responsive on mobile."
```

##### 7. Creating a Sub-task
Link to a parent task using `-p` or `--parent`.

```bash
# First, create the parent task (e.g., it gets ID T15)
backlog create "Implement User Authentication"

# Now, create a sub-task
backlog create "Add Google OAuth login" -p "T15"
```

##### 8. Setting Task Dependencies
Use `--deps` to specify prerequisites.

```bash
# Single dependency
backlog create "Deploy user authentication" --deps "T15"

# Multiple dependencies
backlog create "Integration testing" --deps "T15,T18,T20"
```
{{< /tab >}}
{{< tab "Complex" >}}
##### 9. Complex Example
Combine multiple flags to create a detailed task.

```bash
backlog create "Build new reporting feature" \
  -d "Generate and export monthly PDF reports." \
  -a "drew" \
  -l "feature,frontend,backend" \
  --priority "high" \
  --ac "Report logic is accurate." \
  --ac "User can select a date range." \
  --ac "PDF has correct branding." \
  -p "T23"
```
{{< /tab >}}
{{< /tabpane >}}
{{< /details >}}

### Options

| Flag | Description |
| --- | --- |
| `--ac` | Acceptance criterion (can be specified multiple times). |
| `-a, --assigned` | Assignee for the task (can be specified multiple times). |
| `--deps` | Add a dependency (can be used multiple times). |
| `-d, --description` | Description of the task. |
| `-h, --help` | Help for create. |
| `-l, --labels` | Comma-separated labels for the task. |
| `--notes` | Additional notes for the task. |
| `-p, --parent` | Parent task ID. |
| `--plan` | Implementation plan for the task. |
| `--priority` | Priority of the task (low, medium, high, critical). Default is `medium`. |

{{< details "Inherited Options" >}}
| Flag | Default | Description |
| --- | --- | --- |
| `--auto-commit` | `true` | Auto-committing changes to git repository. |
| `--folder` | `.backlog` | Directory for backlog tasks. |
{{< /details >}}

### SEE ALSO

- [**`backlog`**]({{< relref "backlog" >}}) - The main backlog command.