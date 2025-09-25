## GitHub CLI (`gh`) for AI Agents

This document provides a guide for an AI agent to understand and use the GitHub Command Line Interface (`gh`). `gh` allows you to interact with GitHub from the command line, enabling automation and scripting of various GitHub workflows.

-----

### 1\. Introduction to `gh`

`gh` is GitHub's official command-line tool. It brings pull requests, issues, and other GitHub concepts to the terminal, so you can do all your work in one place. It's a powerful tool for scripting and automation.

**Core Commands:**

  * `gh auth`: Authenticate with your GitHub account.
  * `gh repo`: Create, clone, fork, and view repositories.
  * `gh issue`: Create, list, and manage issues.
  * `gh pr`: Create, list, and manage pull requests.
  * `gh run`: Work with GitHub Actions workflow runs.
  * `gh gist`: Create and manage gists.
  * `gh alias`: Create shortcuts for `gh` commands.
  * `gh api`: Make authenticated requests to the GitHub API.
  * `gh codespace`: Work with GitHub Codespaces.

-----

### 2\. Authentication

Before you can use `gh`, you need to authenticate with your GitHub account. The `gh auth` command group provides several authentication-related functionalities.

**Login:**

The primary way to authenticate is using `gh auth login`. This command will guide you through the authentication process, which usually involves opening a web browser.

```bash
# Start the interactive login process
gh auth login
```

You can also authenticate using a token with the `--with-token` flag.

```bash
# Authenticate with a token from a file
gh auth login --with-token < mytoken.txt
```

**Status:**

To check your authentication status, use `gh auth status`.

```bash
# Check authentication status
gh auth status
```

-----

### 3\. Repository Management (`gh repo`)

The `gh repo` command allows you to manage your repositories.

**Creating a Repository:**

You can create a new repository with `gh repo create`.

```bash
# Create a new public repository and clone it
gh repo create my-new-repo --public --clone

# Create a repository from an existing local directory
gh repo create my-project --private --source=. --remote=upstream
```

**Cloning a Repository:**

To clone a repository, use `gh repo clone`.

```bash
# Clone a repository
gh repo clone cli/cli
```

**Viewing a Repository:**

You can view a repository's README and other information with `gh repo view`.

```bash
# View the current repository
gh repo view

# View a specific repository in the web browser
gh repo view cli/cli --web
```

-----

### 4\. Working with Issues (`gh issue`)

The `gh issue` command group lets you manage issues in your repositories.

**Creating an Issue:**

Create a new issue with `gh issue create`.

```bash
# Create an issue with a title and body
gh issue create --title "Bug found" --body "Something is not working as expected."

# Create an issue with labels and assignees
gh issue create --title "Feature request" --label "enhancement" --assignee "@me"
```

**Listing Issues:**

List issues in a repository with `gh issue list`.

```bash
# List open issues
gh issue list

# List issues with specific labels
gh issue list --label "bug" --label "help wanted"
```

**Closing and Reopening Issues:**

You can close and reopen issues using `gh issue close` and `gh issue reopen`.

```bash
# Close an issue
gh issue close 123

# Reopen an issue with a comment
gh issue reopen 123 --comment "This issue is not resolved."
```

-----

### 5\. Working with Pull Requests (`gh pr`)

The `gh pr` command group is for managing pull requests.

**Creating a Pull Request:**

Create a pull request with `gh pr create`.

```bash
# Create a pull request with a title and body
gh pr create --title "Fixing a bug" --body "This PR fixes a critical bug."

# Create a pull request and request reviews
gh pr create --title "New feature" --body "Adds a new feature." --reviewer monalisa,hubot
```

**Listing Pull Requests:**

List pull requests with `gh pr list`.

```bash
# List open pull requests
gh pr list

# List pull requests by a specific author
gh pr list --author "@me"
```

**Checking Out a Pull Request:**

You can check out a pull request's branch locally with `gh pr checkout`.

```bash
# Check out pull request number 123
gh pr checkout 123
```

**Merging a Pull Request:**

Merge a pull request with `gh pr merge`.

```bash
# Merge a pull request
gh pr merge 123

# Squash and merge a pull request, deleting the branch
gh pr merge 123 --squash --delete-branch
```

-----

### 6\. GitHub Actions (`gh run` and `gh workflow`)

You can manage your GitHub Actions workflows and runs with `gh`.

**Listing Workflow Runs:**

List recent workflow runs with `gh run list`.

```bash
# List recent workflow runs
gh run list

# List workflow runs for a specific branch
gh run list --branch main
```

**Viewing a Workflow Run:**

View the details of a specific run with `gh run view`.

```bash
# View a specific workflow run
gh run view 123456

# View the log for a specific job in a run
gh run view --log --job 789012
```

**Running a Workflow:**

You can manually trigger a workflow with `gh workflow run`.

```bash
# Run the 'build.yml' workflow
gh workflow run build.yml
```

-----

### 7\. Gists (`gh gist`)

`gh` also supports managing GitHub Gists.

**Creating a Gist:**

Create a gist with `gh gist create`.

```bash
# Create a public gist from a file
gh gist create --public my-file.txt

# Create a secret gist with a description
gh gist create my-file.txt -d "A secret gist."
```

**Listing Gists:**

List your gists with `gh gist list`.

```bash
# List all your gists
gh gist list
```

-----

### 8\. Advanced Usage

**GitHub API:**

You can make authenticated requests to the GitHub API using `gh api`.

```bash
# Get information about the authenticated user
gh api user

# List releases for a repository
gh api repos/{owner}/{repo}/releases
```

**Aliases:**

Create aliases for your favorite commands with `gh alias set`.

```bash
# Create an alias 'co' for 'pr checkout'
gh alias set co 'pr checkout'

# Use the alias
gh co 123
```

**Extensions:**

You can extend `gh` with custom commands using extensions.

```bash
# Install an extension
gh extension install owner/gh-extension

# Use the extension
gh extension-command
```

For more detailed information on any command, you can always use `gh help <command>`.
