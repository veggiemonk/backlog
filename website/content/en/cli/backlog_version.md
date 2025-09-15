---
title: "backlog version"
description: "Print the version information for the backlog CLI."
weight: 9
---

## backlog version

Print the version information.

### Synopsis

Prints detailed version information for the `backlog` CLI, including the Git revision, version number, and build time. This is useful for bug reports and checking your installation.

```
backlog version [flags]
```

### Example

```bash
backlog version
```

Example output:
```
Backlog version:
Revision: 7c989dabd2c61a063a23788c18eb39eca408f6a7
Version: v0.0.2-0.20250907193624-7c989dabd2c6
BuildTime: 2025-09-07T19:36:24Z
Dirty: false
```

### Options

| Flag | Description |
| --- | --- |
| `-h, --help` | Help for version. |

{{< details "Inherited Options" >}}
| Flag | Default | Description |
| --- | --- | --- |
| `--auto-commit` | `true` | Auto-committing changes to git repository. |
| `--folder` | `.backlog` | Directory for backlog tasks. |
{{< /details >}}

### SEE ALSO

- [**`backlog`**]({{< relref "backlog" >}}) - The main backlog command.