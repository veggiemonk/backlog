---
layout: page
title: About
nav_order: 6
---

# About Backlog

## Project Background

The context window deteriorates rapidly on large-scale projects. A workaround many developers have found is to ask AI agents to "make a plan for X, write it to a markdown file and keep it updated with the ongoing tasks". This technique has worked incredibly well, making refactoring and other significant code changes more resilient to failures, retries, or rate-limiting from AI models.

However, this approach often cluttered repository root directories with files like `plan.md` or `refactor_this_big_codebase.md`. To solve this problem, Backlog was created as an MCP server that developers could understand and trust to handle these details, providing a much better experience when using AI tools for complex tasks.

## Philosophy

**Simplicity Over Features**: Backlog aims to handle fewer use cases well rather than many use cases poorly. It focuses on the core workflow of task management for developer-AI collaboration.

**AI-First Design**: Unlike traditional task managers designed for humans first, Backlog was designed from the ground up to work seamlessly with AI agents while still being useful for human developers.

**Transparency**: Everything is stored as human-readable Markdown files in Git. No databases, no proprietary formats, no lock-in.

**Developer Control**: You own your data completely. Everything lives in your Git repository and can be edited manually if needed.

## Development Story

While this codebase is mostly hand-written, the documentation, comments, examples, and some tests were generated with AI assistance. The author found that AI agents yield better results when:

1. The project structure is already established exactly as desired
2. The libraries and dependencies are already chosen
3. The code style and patterns are already defined

After establishing these foundations, AI assistance becomes much more effective for implementation details.

## Tools Used

The project was developed using several AI tools (all paid tier):

- **gemini-2.5-pro** - Primary development assistance
- **claude-sonnet4** - Code review and documentation
- **github-copilot** - Code completion and suggestions
- **amp** - Additional AI assistance

## Inspiration

This project is inspired by:

- [Backlog.md](https://github.com/MrLesk/Backlog.md) - Simple markdown-based task tracking
- [TaskWing](https://github.com/josephgoksu/TaskWing) - CLI task management

But Backlog aims to be simpler while being specifically optimized for AI-friendly workflows.

For MCP server implementation, excellent examples can be found in the [go-sdk](https://github.com/modelcontextprotocol/go-sdk).

For production MCP server use cases, [GoogleCloudPlatform/gke-mcp](https://github.com/GoogleCloudPlatform/gke-mcp) is highly recommended.

## Example Usage

You can see Backlog in action by examining the [.backlog](./.backlog) folder in this repository, which contains the actual tasks used to develop this project.

The backlog was generated using this prompt:

```
If you were to recreate this project from scratch, make a plan and break it down into tasks using backlog break down prompt.
Write that plan to a markdown file called "./docs/plan.md".
Check that the plan in ./docs/plan.md is consistent with the list of tasks in the backlog.
Add implementation plan to relevant tasks
All tasks should have at least one acceptance criteria.

Read the full instructions for backlog: ./internal/mcp/prompt.md
```

## Contributing

Backlog is open source and welcomes contributions. See the [Development Guide](development/) for technical details and the [GitHub repository](https://github.com/veggiemonk/backlog) for the latest code.

## License

This project is released under the MIT License. See the LICENSE file in the repository for details.

## Contact

- **Repository**: [github.com/veggiemonk/backlog](https://github.com/veggiemonk/backlog)
- **Issues**: [GitHub Issues](https://github.com/veggiemonk/backlog/issues)
- **Author**: veggiemonk

---

*"The goal is to provide a frictionless collaboration between AI agents and developers"*