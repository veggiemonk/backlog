# Contributing to Backlog

First off, thank you for considering contributing to Backlog! It's people like you that make open source software great. We welcome any and all contributions, from bug reports to feature requests and code contributions.

## How to Contribute

There are many ways to contribute to the project:

-   **Reporting Bugs**: If you find a bug, please open an issue and provide as much detail as possible, including steps to reproduce it.
-   **Suggesting Enhancements**: If you have an idea for a new feature or an improvement to an existing one, please open an issue to discuss it.
-   **Writing Code**: If you want to contribute code, please follow the steps below.

## Development Workflow

1.  **Fork the Repository**: Start by forking the repository to your own GitHub account.

2.  **Clone the Repository**: Clone your fork to your local machine:

    ```sh
    git clone https://github.com/YOUR_USERNAME/backlog.git
    cd backlog
    ```

3.  **Create a Branch**: Create a new branch for your changes:

    ```sh
    git checkout -b my-feature-branch
    ```

4.  **Set Up the Development Environment**: Make sure you have Go installed. Then, install the necessary tools:

    ```sh
    make setup
    ```

5.  **Make Your Changes**: Now you can start making your changes to the code.

6.  **Follow the Code Style**: Be sure to follow the code style guidelines outlined in the [`docs/code-style.md`](./docs/code-style.md) document.

7.  **Run Tests**: Before submitting your changes, make sure that all tests pass:

    ```sh
    make test
    ```

8.  **Lint Your Code**: Ensure your code is properly formatted and linted:

    ```sh
    make lint
    ```

9.  **Commit Your Changes**: Commit your changes with a clear and descriptive commit message:

    ```sh
    git commit -m "feat: Add new feature"
    ```

10. **Push to Your Fork**: Push your changes to your fork:

    ```sh
    git push origin my-feature-branch
    ```

11. **Submit a Pull Request**: Open a pull request from your fork to the main `backlog` repository. In the pull request description, please explain the changes you made and link to any relevant issues.

## Pull Request Process

1.  A maintainer will review your pull request and may suggest some changes.
2.  Once the pull request is approved, it will be merged into the main branch.

Thank you for your contribution!
