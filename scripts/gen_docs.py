import mkdocs_gen_files
import re


def get_section(markdown_content, section_title):
    """Extracts a section from markdown content."""
    # This pattern looks for a ## heading and captures everything until the next ## heading.
    pattern = re.compile(
        r"(^## " + re.escape(section_title) + r".*?)(?=^## )", re.MULTILINE | re.DOTALL
    )
    match = pattern.search(markdown_content)
    if match:
        return match.group(1).strip()

    # This is a fallback for the case where the section is the last one in the file.
    pattern = re.compile(
        r"(^## " + re.escape(section_title) + r".*)", re.MULTILINE | re.DOTALL
    )
    match = pattern.search(markdown_content)
    if match:
        return match.group(1).strip()

    return None


with open("README.md", "r") as f:
    readme = f.read()

# Create a page for the prompt instructions
with open("internal/mcp/prompt.md", "r") as f:
    prompt_content = f.read()

with mkdocs_gen_files.open("prompts/mcp.md", "w") as f:
    f.write("---\n")
    f.write("layout: page\n")
    f.write("title: Prompt to use Backlog CLI\n")
    f.write("nav_order: 2\n")
    f.write("---\n\n")
    f.write(prompt_content)

def replace_internal_links(content):
    """Replaces internal links with their corresponding URLs for the docs site."""
    return (
        content.replace("(./internal/mcp/prompt.md)", "(prompts/mcp.md)")
        .replace(
            "(./.backlog)",
            "(https://github.com/veggiemonk/backlog/tree/main/.backlog)",
        )
        .replace(
            "(./.gemini)",
            "(https://github.com/veggiemonk/backlog/tree/main/.gemini)",
        )
        .replace(
            "(./.claude)",
            "(https://github.com/veggiemonk/backlog/tree/main/.claude)",
        )
    )


# Fix the link in the README content for the docs site
readme_for_docs = replace_internal_links(readme)

# Generate index.md from the full README
with mkdocs_gen_files.open("index.md", "w") as f:
    f.write("---\n")
    f.write("layout: home\n")
    f.write("title: Home\n")
    f.write("nav_order: 1\n")
    f.write("---\n\n")
    f.write(readme_for_docs)


# Define the sections to extract and the files to generate
sections = {
    "Installation": ("installation.md", 2),
    "Quick Start": ("getting-started.md", 3),
    "Usage Examples": ("usage-examples.md", 4),
    "AI Agent Integration": ("ai-integration.md", 5),
}

for title, (filename, nav_order) in sections.items():
    content = get_section(readme, title)
    if content:
        content_for_docs = replace_internal_links(content)
        with mkdocs_gen_files.open(filename, "w") as f:
            f.write("---\n")
            f.write("layout: page\n")
            f.write(f"title: {title}\n")
            f.write(f"nav_order: {nav_order}\n")
            f.write("---\n\n")
            f.write(content_for_docs)
