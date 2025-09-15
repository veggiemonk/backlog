---
title: "About"
description: "Project background, philosophy, and inspiration"
weight: 60
---

# About Backlog

## Project Background

{{< alert >}}
The context window deteriorates rapidly on large-scale projects. A workaround many developers have found is to ask AI agents to "make a plan for X, write it to a markdown file and keep it updated with the ongoing tasks".
{{< /alert >}}

This technique has worked incredibly well, making refactoring and other significant code changes more resilient to failures, retries, or rate-limiting from AI models.

However, this approach often cluttered repository root directories with files like `plan.md` or `refactor_this_big_codebase.md`. To solve this problem, **Backlog** was created as an MCP server that developers could understand and trust to handle these details, providing a much better experience when using AI tools for complex tasks.

## Philosophy

**Simplicity Over Features**
: Backlog aims to handle fewer use cases well rather than many use cases poorly. It focuses on the core workflow of task management for developer-AI collaboration.

**AI-First Design**
: Unlike traditional task managers designed for humans first, Backlog was designed from the ground up to work seamlessly with AI agents while still being useful for human developers.

**Transparency**
: Everything is stored as human-readable Markdown files in Git. No databases, no proprietary formats, no lock-in.

**Developer Control**
: You own your data completely. Everything lives in your Git repository and can be edited manually if needed.

## Development Story

{{< alert >}}
While this codebase is mostly hand-written, the documentation, comments, examples, and some tests were generated with AI assistance.
{{< /alert >}}

The author found that AI agents yield better results when:

1. **The project structure is already established** exactly as desired
2. **The libraries and dependencies are already chosen** and configured
3. **The code style and patterns are already defined** consistently

After establishing these foundations, AI assistance becomes much more effective for implementation details.

## Tools Used

The project was developed using several AI tools (all paid tier):

**Primary Tools:**
- **gemini-2.5-pro** - Primary development assistance
- **claude-sonnet4** - Code review and documentation

**Supporting Tools:**
- **github-copilot** - Code completion and suggestions
- **amp** - Additional AI assistance

## Inspiration

This project is inspired by:

{{< details "Backlog.md" >}}
[Backlog.md](https://github.com/MrLesk/Backlog.md) - Simple markdown-based task tracking that demonstrated the power of storing tasks as plain markdown files.
{{< /details >}}

{{< details "TaskWing" >}}
[TaskWing](https://github.com/josephgoksu/TaskWing) - CLI task management that showed how effective command-line interfaces can be for developer workflows.
{{< /details >}}

But Backlog aims to be simpler while being specifically optimized for AI-friendly workflows.

### MCP Server Resources

For MCP server implementation, excellent examples can be found in:

- [**go-sdk**](https://github.com/modelcontextprotocol/go-sdk) - Official Go SDK for MCP servers
- [**GoogleCloudPlatform/gke-mcp**](https://github.com/GoogleCloudPlatform/gke-mcp) - Production MCP server example

## Example Usage

{{< alert >}}
You can see Backlog in action by examining the [.backlog](https://github.com/veggiemonk/backlog/tree/main/.backlog) folder in this repository, which contains the actual tasks used to develop this project.
{{< /alert >}}

The backlog was generated using this prompt:

{{< details "Original Planning Prompt" >}}
```
If you were to recreate this project from scratch, make a plan and break it down into tasks using backlog break down prompt.
Write that plan to a markdown file called "./docs/plan.md".
Check that the plan in ./docs/plan.md is consistent with the list of tasks in the backlog.
Add implementation plan to relevant tasks
All tasks should have at least one acceptance criteria.

Read the full instructions for backlog: ./internal/mcp/prompt.md
```
{{< /details >}}

## Contributing

Backlog is open source and welcomes contributions. See the [Development Guide]({{< relref "development" >}}) for technical details.

**Resources:**
- **Repository**: [github.com/veggiemonk/backlog](https://github.com/veggiemonk/backlog)
- **Issues**: [GitHub Issues](https://github.com/veggiemonk/backlog/issues)
- **Discussions**: [GitHub Discussions](https://github.com/veggiemonk/backlog/discussions)

**Quick Links:**
- [**Getting Started**]({{< relref "getting-started" >}}) - Installation and basic usage
- [**CLI Reference**]({{< relref "cli" >}}) - Complete command documentation
- [**AI Integration**]({{< relref "ai-integration" >}}) - MCP setup guide

## License

This project is released under the **MIT License**. See the [LICENSE](https://github.com/veggiemonk/backlog/blob/main/LICENSE) file in the repository for details.

## Contact

**Project Details:**
- **Repository**: [github.com/veggiemonk/backlog](https://github.com/veggiemonk/backlog)
- **Issues**: [GitHub Issues](https://github.com/veggiemonk/backlog/issues)
- **Author**: [veggiemonk](https://github.com/veggiemonk)

**Documentation:**
- **User Guide**: [Getting Started]({{< relref "getting-started" >}})
- **Technical Docs**: [Development]({{< relref "development" >}})
- **API Reference**: [CLI Commands]({{< relref "cli" >}})

---

{{< alert >}}
*"The goal is to provide a frictionless collaboration between AI agents and developers"*
{{< /alert >}}