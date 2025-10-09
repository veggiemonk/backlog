---
layout: page
title: backlog instructions
---

# NAME

instructions - instructions for agents to learn to use backlog

# SYNOPSIS

instructions

```
[--mode]=[value]
```

# DESCRIPTION

Instructions for agents to learn to use backlog by including them into a prompt.

Examples:

backlog instructions               # outputs the instructions for agents to use the cli.
backlog instructions --mode cli    # outputs the instructions for agents to use the cli.
backlog instructions --mode mcp    # outputs the instructions for agents to use MCP.
backlog instructions >> AGENTS.md  # add instructions to agent base prompt.


**Usage**:

```
instructions [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--mode**="": which mode the agent will use backlog: (cli|mcp) (default: cli)

