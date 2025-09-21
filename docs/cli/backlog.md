## backlog

backlog is a git-native, markdown-based task manager

### Synopsis

A Git-native, Markdown-based task manager for developers and AI agents.
Backlog helps you manage tasks within your git repository.

```
backlog [flags]
```

### Options

```
      --auto-commit         Auto-committing changes to git repository
      --folder string       Directory for backlog tasks (default ".backlog")
  -h, --help                help for backlog
      --log-file string     Log file path (defaults to stderr)
      --log-format string   Log format (json, text) (default "text")
      --log-level string    Log level (debug, info, warn, error) (default "info")
      --max-limit int       Maximum limit for pagination (default 1000)
      --page-size int       Default page size for pagination (default 25)
```

### SEE ALSO

* [backlog archive](backlog_archive.md)	 - Archive a task
* [backlog create](backlog_create.md)	 - Create a new task
* [backlog edit](backlog_edit.md)	 - Edit an existing task
* [backlog instructions](backlog_instructions.md)	 - instructions for agents to learn to use backlog
* [backlog list](backlog_list.md)	 - List all tasks
* [backlog mcp](backlog_mcp.md)	 - Start the MCP server
* [backlog search](backlog_search.md)	 - Search tasks by content
* [backlog version](backlog_version.md)	 - Print the version information
* [backlog view](backlog_view.md)	 - View a task by providing its ID

