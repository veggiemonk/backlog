---
title: "backlog archive"
description: "Archive a task, moving it to the archived directory."
weight: 7
---

## backlog archive

Archive a task.

### Synopsis

Archives a task, moving it to the archived directory and marking its status as `archived`. This is the best way to clean up your board without losing history.

```
backlog archive <task-id> [flags]
```

{{< hint type=tip >}}
You can use a partial ID to identify the task, as long as it's unique. For example, `backlog archive T01` is equivalent to `backlog archive T01-some_title.md`.
{{< /hint >}}

### Options

| Flag | Description |
| --- | --- |
| `-h, --help` | Help for archive. |

{{< expand "Inherited Options" >}}
| Flag | Default | Description |
| --- | --- | --- |
| `--auto-commit` | `true` | Auto-committing changes to git repository. |
| `--folder` | `.backlog` | Directory for backlog tasks. |
{{< /expand >}}

### SEE ALSO

- [**`backlog`**]({{< relref "backlog" >}}) - The main backlog command.