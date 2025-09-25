# A Guide to `tmux` for Humans and AI Agents

This document provides a guide to understanding and using `tmux`, a terminal multiplexer. It is designed to be read by humans and parsed by AI agents to learn how to manage terminal sessions effectively.

## 1. What is `tmux`?

`tmux` is a **terminal multiplexer**. It allows you to create, manage, and switch between multiple virtual terminals from a single terminal window. This enables:

1.  **Session Persistence**: Keep your sessions and running processes alive even if you get disconnected from the server.
2.  **Multitasking**: View and work with multiple command-line programs in a single screen using "panes" (split-screen) and "windows" (tabs).
3.  **Session Sharing**: Allow multiple users to connect to the same terminal session for pair programming or collaboration.

## 2. Core Concepts

`tmux` has a simple hierarchical structure:

-   **Server**: A background process that manages all sessions. It starts automatically when you run your first `tmux` command.
-   **Session**: A single workspace that contains one or more windows. You can have multiple sessions for different projects (e.g., `work`, `personal`).
-   **Window**: A single screen within a session, similar to a browser tab. Each window fills the entire terminal screen.
-   **Pane**: A rectangular section of a window where a separate shell or program is running. You can split a window into multiple panes.

### The Prefix Key

All `tmux` keyboard shortcuts are triggered by first pressing a **prefix key**. The default prefix is `Ctrl-b`. After pressing the prefix, you press the command key (e.g., `c` to create a new window).

To send an actual `Ctrl-b` keystroke to the application running in a pane, you press `Ctrl-b` twice.

## 3. Common Use Cases & Commands

### Use Case 1: Session Persistence

This is the most powerful feature of `tmux`. You can start a long-running process (like a build, server, or training job), detach from the session, log off, and re-attach later to find it still running.

#### Commands

-   **Start a new named session:**
    ```sh
    tmux new -s my_session
    ```

-   **List all running sessions:**
    ```sh
    tmux ls
    ```

-   **Detach from the current session (leaves it running):**
    -   Press `Ctrl-b` then `d`.

-   **Attach to the most recent session:**
    ```sh
    tmux attach
    ```

-   **Attach to a specific named session:**
    ```sh
    tmux attach -t my_session
    ```

-   **Kill a specific session:**
    ```sh
    tmux kill-session -t my_session
    ```

### Use Case 2: Multitasking with Windows and Panes

You can organize your workspace for a project with multiple related command-line tasks.

#### Managing Windows (Tabs)

| Action | Keystroke | Description |
| :--- | :--- | :--- |
| Create a new window | `Ctrl-b` `c` | A new window is created with a shell prompt. |
| Go to the next window | `Ctrl-b` `n` | Cycles forward through your windows. |
| Go to the previous window | `Ctrl-b` `p` | Cycles backward through your windows. |
| Go to a specific window | `Ctrl-b` `[0-9]`| Jumps directly to the window number. |
| Rename the current window | `Ctrl-b` `,` | Useful for organizing your work. |

#### Managing Panes (Splits)

| Action | Keystroke | Description |
| :--- | :--- | :--- |
| Split pane vertically | `Ctrl-b` `%` | Creates a new pane to the right. |
| Split pane horizontally | `Ctrl-b` `"` | Creates a new pane below. |
| Navigate between panes | `Ctrl-b` `[arrow key]` | Moves the cursor to the pane in that direction. |
| Close the current pane | `Ctrl-b` `x` | You will be prompted to confirm. |
| Toggle pane zoom | `Ctrl-b` `z` | Expands the current pane to fill the window. Press again to unzoom. |

### Use Case 3: Pair Programming

Multiple users can connect to the same `tmux` session to share a terminal.

1.  **User 1**: Starts a new session on a shared server.
    ```sh
    tmux new -s pair_programming
    ```
2.  **User 2**: Attaches to that same session.
    ```sh
    tmux attach -t pair_programming
    ```

Both users will see the same screen and can type in the same terminal. `tmux` will resize the view to the smallest user's terminal.

## 4. Basic Configuration (`~/.tmux.conf`)

You can customize `tmux` by creating a configuration file at `~/.tmux.conf`.

**Example `~/.tmux.conf`:**

```
# Change the prefix key to Ctrl-a (like the 'screen' command)
set-option -g prefix C-a
unbind-key C-b
bind-key C-a send-prefix

# Enable mouse mode for scrolling and pane selection
set -g mouse on

# Use more intuitive split commands
bind | split-window -h
bind - split-window -v
unbind '''%'''
unbind '''"'''
```

After editing the file, you need to reload the configuration from within a `tmux` session:
-   Press `Ctrl-b` then `:` to open the command prompt.
-   Type `source-file ~/.tmux.conf` and press Enter.

## 5. Command Cheatsheet

| Action              | Command                       | Keystroke           |
| :------------------ | :---------------------------- | :------------------ |
| **Sessions**        |                               |                     |
| New Session         | `tmux new -s <name>`          |                     |
| List Sessions       | `tmux ls`                     |                     |
| Attach to Session   | `tmux attach -t <name>`       |                     |
| Detach from Session |                               | `Ctrl-b` `d`        |
| Kill Session        | `tmux kill-session -t <name>` |                     |
| **Windows**         |                               |                     |
| New Window          |                               | `Ctrl-b` `c`        |
| Next Window         |                               | `Ctrl-b` `n`        |
| Previous Window     |                               | `Ctrl-b` `p`        |
| Rename Window       |                               | `Ctrl-b` `,`        |
| **Panes**           |                               |                     |
| Split Vertically    |                               | `Ctrl-b` `%`        |
| Split Horizontally  |                               | `Ctrl-b` `"`        |
| Navigate Panes      |                               | `Ctrl-b` `[arrows]` |
| Close Pane          |                               | `Ctrl-b` `x`        |
| Zoom/Unzoom Pane    |                               | `Ctrl-b` `z`        |

