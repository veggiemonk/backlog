---
layout: page
title: backlog instructions
---

## backlog instructions

instructions for agents to learn to use backlog

### Synopsis

Instructions for agents to learn to use backlog by including them into a prompt.

```
backlog instructions [flags]
```

### Examples

```

backlog instructions               # outputs the instructions for agents to use the cli.
backlog instructions --mode cli    # outputs the instructions for agents to use the cli.
backlog instructions --mode mcp    # outputs the instructions for agents to use MCP.
backlog instructions >> AGENTS.md  # add instructions to agent base prompt.

```

### Options

```
  -h, --help          help for instructions
      --mode string   which mode the agent will use backlog: (cli|mcp) (default "cli")
```

### Options inherited from parent commands

```
      --auto-commit         Auto-committing changes to git repository (default true)
      --folder string       Directory for backlog tasks (default ".backlog")
      --log-file string     Log file path (defaults to stderr)
      --log-format string   Log format (json, text) (default "text")
      --log-level string    Log level (debug, info, warn, error) (default "info")
```

### SEE ALSO

* [backlog](backlog.md)	 - Backlog is a git-native, markdown-based task manager

