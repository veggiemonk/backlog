---
layout: page
title: backlog view
---

# NAME

view - View a task by providing its ID

# SYNOPSIS

view

```
[--json|-j]
```

# DESCRIPTION


View a task by providing its ID. You can output in markdown or JSON format.

Examples: 
```

  backlog view T01           # View task T01 in markdown format
  backlog view T01 --json    # View task T01 in JSON format
  backlog view T01 -j        # View task T01 in JSON format (short flag)

```

**Usage**:

```
backlog view <id>
```

# GLOBAL OPTIONS

**--json, -j**: Print JSON output

