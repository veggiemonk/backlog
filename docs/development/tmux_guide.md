# A Guide to `tmux` for AI Agents

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

### AI Agent Interaction

For AI agents, the concept of a "prefix key" and keyboard shortcuts is irrelevant. AI agents interact with `tmux` directly through shell commands. All operations, such as creating windows, splitting panes, or sending commands to a pane, are executed using `tmux` commands in the shell, as detailed in the "Programmatic Control for Agents" section. This direct command-line interface allows agents to orchestrate `tmux` sessions programmatically without needing to simulate keyboard input.

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
    ```sh
    tmux detach
    ```

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

| Action | Command | Description |
| :--- | :--- | :--- |
| Create a new window | `tmux new-window` | A new window is created with a shell prompt. |
| Go to the next window | `tmux next-window` | Cycles forward through your windows. |
| Go to the previous window | `tmux previous-window` | Cycles backward through your windows. |
| Go to a specific window | `tmux select-window -t <window-number>`| Jumps directly to the window number. |
| Rename the current window | `tmux rename-window <new-name>` | Useful for organizing your work. |

#### Managing Panes (Splits)

| Action | Command | Description |
| :--- | :--- | :--- |
| Split pane vertically | `tmux split-window -h` | Creates a new pane to the right. |
| Split pane horizontally | `tmux split-window -v` | Creates a new pane below. |
| Navigate between panes | `tmux select-pane -U/D/L/R` | Moves the cursor to the pane in that direction. |
| Close the current pane | `tmux kill-pane` | You will be prompted to confirm. |
| Toggle pane zoom | `tmux resize-pane -Z` | Expands the current pane to fill the window. Press again to unzoom. |

## 3. Programmatic Control for Agents

For an AI agent to orchestrate complex tasks, it needs to use shell commands instead of key bindings. The following commands are essential for programmatic control.

### Targeting
Commands often require a target, specified with the `-t` flag. The target syntax is `session_name:window_index.pane_index`.
- `my_session:0.1` refers to the second pane of the first window in the session named `my_session`.

### Capturing Output
To see the output of a command, an agent must capture the contents of a pane.

### Checking for Existence
To avoid errors, an agent can check if a session already exists before trying to create it.

| Action | Command | Description |
| :--- | :--- | :--- |
| **Send Keystrokes** | `tmux send-keys -t <target> "command" C-m` | Sends a command string and an "Enter" (`C-m`) to the target pane. This is the primary way to run commands in an existing pane. |
| **New Window** | `tmux new-window -t <target-session> -n <name>` | Creates a new window in the target session and gives it a name. |
| **Split Pane** | `tmux split-window -h -t <target-pane>` | Splits a target pane horizontally (creating a vertical line). Use `-v` for a vertical split. |
| **Select Pane** | `tmux select-pane -t <target-pane>` | Selects a specific pane (makes it the active one). |
| **List Panes** | `tmux list-panes -s -t <target-session>` | Lists all panes in a session with their details, allowing an agent to discover the layout and pane indexes. |
| **Capture Pane** | `tmux capture-pane -p -t <target-pane>` | Captures the full contents (including scrollback history) of a target pane and prints it to standard output. |
| **Check for Session** | `tmux has-session -t <session-name>` | Returns a zero exit code if the session exists, and a non-zero exit code otherwise. This is useful for scripting. |


## 4. Example Workflow for Agents

Here is a step-by-step example of how an AI agent can use `tmux` to run a process in the background, check its output, and then terminate it.

1.  **Check if a session named `my-app` already exists.**
    ```sh
    tmux has-session -t my-app
    ```

2.  **If it doesn't exist, create it.**
    ```sh
    tmux new-session -d -s my-app
    ```
    *   The `-d` flag starts the session in detached mode, so it doesn't attach to it.

3.  **Create a new window for the application.**
    ```sh
    tmux new-window -t my-app -n server
    ```

4.  **Split the window into two panes.**
    ```sh
    tmux split-window -h -t my-app:server
    ```

5.  **In the first pane (pane 0), start a web server.**
    ```sh
    tmux send-keys -t my-app:server.0 "python3 -m http.server 8080" C-m
    ```

6.  **In the second pane (pane 1), use `curl` to check if the server is running.**
    ```sh
    tmux send-keys -t my-app:server.1 "curl -s http://localhost:8080" C-m
    ```

7.  **Capture the output of the `curl` command from pane 1.**
    ```sh
    tmux capture-pane -p -t my-app:server.1
    ```
    *   The agent can then parse this output to verify the server's response.

8.  **Kill the session to clean up.**
    ```sh
    tmux kill-session -t my-app
    ```

## 5. Command Cheatsheet

| Action | Command |
| :--- | :--- |
| **Sessions** | |
| New Session | `tmux new -s <name>` |
| List Sessions | `tmux ls` |
| Attach to Session | `tmux attach -t <name>` |
| Detach from Session | `tmux detach` |
| Kill Session | `tmux kill-session -t <name>` |
| **Windows** | |
| New Window | `tmux new-window -n <name>` |
| Next Window | `tmux next-window` |
| Previous Window | `tmux previous-window` |
| Rename Window | `tmux rename-window <new-name>` |
| **Panes** | |
| Split Vertically | `tmux split-window -v` |
| Split Horizontally | `tmux split-window -h` |
| Navigate Panes | `tmux select-pane -U/D/L/R` |
| Close Pane | `tmux kill-pane` |
| Zoom/Unzoom Pane | `tmux resize-pane -Z` |
| **Agent Control** | |
| Send Command | `tmux send-keys -t <target> "cmd" C-m` |
