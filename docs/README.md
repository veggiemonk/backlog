# Backlog Documentation

This directory contains the Jekyll-based documentation website for the Backlog project.

## Structure

- `index.md` - Homepage
- `getting-started.md` - Installation and basic usage guide
- `ai-integration.md` - MCP server setup for AI agents
- `about.md` - Project background and philosophy
- `cli/` - Complete CLI command reference
- `development/` - Technical documentation for contributors
- `_layouts/` - Jekyll layout templates
- `assets/` - CSS and other static assets

## Local Development

To run the documentation locally:

```bash
# Install dependencies
bundle install

# Serve the site
bundle exec jekyll serve

# Open http://localhost:4000/backlog
```

## Deployment

The documentation is automatically deployed to GitHub Pages via GitHub Actions when changes are pushed to the main branch.

The site will be available at: https://veggiemonk.github.io/backlog