---
title: "CLI Reference"
description: "Complete command-line interface documentation"
weight: 30
geekdocCollapseSection: true
---

# CLI Reference

Complete command-line interface documentation for the backlog tool.

## Overview

Backlog provides a comprehensive CLI for managing tasks with the following commands:

<div class="row">
  <div class="col-md-6">
    <h3>Core Commands</h3>
    <ul>
      <li><a href="backlog"><code>backlog</code></a> - Main command with global options</li>
      <li><a href="backlog_create"><code>backlog create</code></a> - Create new tasks</li>
      <li><a href="backlog_list"><code>backlog list</code></a> - List and filter tasks</li>
      <li><a href="backlog_view"><code>backlog view</code></a> - View task details</li>
      <li><a href="backlog_edit"><code>backlog edit</code></a> - Edit existing tasks</li>
    </ul>
  </div>
  <div class="col-md-6">
    <h3>Additional Commands</h3>
    <ul>
      <li><a href="backlog_search"><code>backlog search</code></a> - Search tasks by content</li>
      <li><a href="backlog_archive"><code>backlog archive</code></a> - Archive completed tasks</li>
      <li><a href="backlog_mcp"><code>backlog mcp</code></a> - Start MCP server for AI agents</li>
      <li><a href="backlog_version"><code>backlog version</code></a> - Show version information</li>
    </ul>
  </div>
</div>

## Quick Examples

### Create Tasks

{{< tabpane >}}
{{< tab "Simple Task" >}}
```bash
backlog create "Fix login bug"
```
{{< /tab >}}
{{< tab "With Metadata" >}}
```bash
backlog create "Implement OAuth" \
  -d "Add Google and GitHub OAuth support" \
  -a "alice,bob" \
  -l "auth,feature" \
  --priority "high" \
  --ac "Google OAuth works" \
  --ac "GitHub OAuth works"
```
{{< /tab >}}
{{< tab "Subtask" >}}
```bash
backlog create "OAuth token validation" -p "T01"
```
{{< /tab >}}
{{< /tabpane >}}

### List and Filter

{{< tabpane >}}
{{< tab "All Tasks" >}}
```bash
backlog list
```
{{< /tab >}}
{{< tab "By Status" >}}
```bash
backlog list --status "in-progress"
```
{{< /tab >}}
{{< tab "By Assignee" >}}
```bash
backlog list --assignee "alice"
```
{{< /tab >}}
{{< tab "Subtasks" >}}
```bash
backlog list --parent "T01"
```
{{< /tab >}}
{{< /tabpane >}}

### View and Edit

<div class="row">
  <div class="col-md-6">
    <h4>View Details</h4>
    <pre><code class="language-bash">backlog view T01</code></pre>
    <h4>Update Status</h4>
    <pre><code class="language-bash">backlog edit T01 --status "done"</code></pre>
  </div>
  <div class="col-md-6">
    <h4>Add Assignee</h4>
    <pre><code class="language-bash">backlog edit T01 --add-assignee "bob"</code></pre>
    <h4>Add Labels</h4>
    <pre><code class="language-bash">backlog edit T01 --add-label "urgent"</code></pre>
  </div>
</div>

### Search and Archive

<div class="row">
  <div class="col-md-6">
    <h4>Search Content</h4>
    <pre><code class="language-bash">backlog search "authentication"</code></pre>
  </div>
  <div class="col-md-6">
    <h4>Archive Task</h4>
    <pre><code class="language-bash">backlog archive T01</code></pre>
  </div>
</div>

### AI Integration

{{< tabpane >}}
{{< tab "STDIO Transport" >}}
```bash
backlog mcp
```
{{< /tab >}}
{{< tab "HTTP Transport" >}}
```bash
backlog mcp --http --port 8106
```
{{< /tab >}}
{{< /tabpane >}}

## Global Options

{{< alert title="Tip" color="info" >}}
All commands support these global options:
{{< /alert >}}

- `--help, -h`: Show help information
- `--version`: Show version information

## Output Formats

Many commands support multiple output formats:

<div class="row">
  <div class="col-md-6">
    <ul>
      <li><strong>Table</strong> (default): Human-readable table format</li>
      <li><strong>JSON</strong>: Machine-readable JSON format</li>
      <li><strong>Markdown</strong>: Markdown format for documentation</li>
    </ul>
  </div>
  <div class="col-md-6">
    <pre><code class="language-bash">backlog list --format json
backlog view T01 --format markdown</code></pre>
  </div>
</div>

## Task ID Format

{{< alert title="Note" >}}
Tasks use hierarchical IDs with dot notation:
{{< /alert >}}

- `T01` - Root task
- `T01.01` - Subtask of T01
- `T01.01.01` - Sub-subtask of T01.01

IDs are automatically assigned when creating tasks, or you can specify a parent with the `-p` flag.

## File Structure

Tasks are stored as Markdown files in the `.backlog/` directory:

```
.backlog/
├── T01-implement_auth.md
├── T01.01-oauth_setup.md
├── T02-frontend_work.md
└── archived/
    └── T03-completed_task.md
```

{{< alert title="Tip" color="info" >}}
Each command in this reference includes detailed usage information, examples, and available options.
{{< /alert >}}