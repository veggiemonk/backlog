package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/cmd"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var (
	dirDocs           = filepath.Join(".", "docs")
	dirRef            = filepath.Join(dirDocs, "reference")
	sectionsToExtract = []string{"usage_examples", "quick_start", "ai_agent_integration"}
)

var dirReplacer = strings.NewReplacer(
	"(./.backlog)", "(https://github.com/veggiemonk/backlog/tree/main/.backlog)",
	"(./.gemini)", "(https://github.com/veggiemonk/backlog/tree/main/.gemini)",
	"(./.claude)", "(https://github.com/veggiemonk/backlog/tree/main/.claude)",
)

var (
	cliPromptPath = "internal/mcp/prompt-cli.md"
	mcpPromptPath = "internal/mcp/prompt-mcp.md"
)

func main() {
	format := flag.String("format", "markdown", "markdown")
	flag.Parse()

	checkErr(func() error { return os.MkdirAll(dirDocs, 0o750) })
	checkErr(func() error { return genReference(dirRef, *format) })
	checkErr(splitReadMe)
	checkErr(addPromptCLI)
	checkErr(addPromptMCP)
	checkErr(addAGENTSmd)
	checkErr(addIndex)
}

func genReference(out, format string) error {
	if format != "markdown" {
		return fmt.Errorf("unknown format: %s", format)
	}

	if err := os.MkdirAll(out, 0o750); err != nil {
		return err
	}

	root := cmd.NewCommand(cmd.WithSkipLogging(true))

	if err := writeMarkdownDoc(out, "backlog.md", "backlog", root); err != nil {
		return err
	}

	for _, sub := range root.Commands {
		if sub.Hidden {
			continue
		}
		filename := fmt.Sprintf("backlog_%s.md", sub.Name)
		title := fmt.Sprintf("backlog %s", sub.Name)
		if err := writeMarkdownDoc(out, filename, title, sub); err != nil {
			return err
		}
	}

	return nil
}

func writeMarkdownDoc(outDir, filename, title string, command *cli.Command) error {
	md, err := docs.ToMarkdown(command)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.WriteString("layout: page\n")
	buf.WriteString(fmt.Sprintf("title: %s\n", title))
	buf.WriteString("---\n\n")
	buf.WriteString(md)

	path := filepath.Join(outDir, filename)
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func addIndex() error {
	b, err := os.ReadFile("README.md")
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.WriteString("layout: home\n")
	buf.WriteString("title: Home\n")
	buf.WriteString("nav_order:  1\n")
	buf.WriteString("---\n\n")
	content := dirReplacer.Replace(string(b))
	content = strings.ReplaceAll(content, "(./"+cliPromptPath+")", "(prompts/cli.md)")
	content = strings.ReplaceAll(content, "(./"+mcpPromptPath+")", "(prompts/mcp.md)")
	buf.WriteString(content)
	path := filepath.Join(dirDocs, "index.md")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}

func addPromptCLI() error {
	b, err := os.ReadFile(cliPromptPath)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.WriteString("layout: page\n")
	buf.WriteString("title: Prompt to use Backlog CLI\n")
	buf.WriteString("nav_order: 1\n")
	buf.WriteString("---\n\n")
	content := dirReplacer.Replace(string(b))
	buf.WriteString(content)
	path := filepath.Join(dirDocs, "prompts", "cli.md")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

func addPromptMCP() error {
	b, err := os.ReadFile(mcpPromptPath)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.WriteString("layout: page\n")
	buf.WriteString("title: Prompt to use Backlog MCP\n")
	buf.WriteString("nav_order: 2\n")
	buf.WriteString("---\n\n")
	content := dirReplacer.Replace(string(b))
	buf.WriteString(content)
	path := filepath.Join(dirDocs, "prompts", "mcp.md")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

func addAGENTSmd() error {
	b, err := os.ReadFile("AGENTS.md")
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.WriteString("layout: page\n")
	buf.WriteString("title: AGENTS.md\n")
	buf.WriteString("nav_order: 3\n")
	buf.WriteString("---\n\n")
	content := dirReplacer.Replace(string(b))
	content = strings.ReplaceAll(content, "(./internal/mcp/prompt-cli.md)", "(../prompts/cli.md)")
	buf.WriteString(content)
	path := filepath.Join(dirDocs, "prompts", "AGENTS.md")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

func splitReadMe() error {
	b, err := os.ReadFile("README.md")
	if err != nil {
		return err
	}
	p := goldmark.DefaultParser()
	doc := p.Parse(text.NewReader(b))
	sections := extractSections(doc, b)
	for _, section := range sections {
		if section.Title == "" {
			continue
		}
		slug := slugify(section.Title)
		log.Println("section", slug)
		if !slices.Contains(sectionsToExtract, slug) {
			continue
		}
		path := filepath.Join(dirDocs, section.Filename)
		if err := os.WriteFile(path, []byte(section.Content), 0o644); err != nil {
			return fmt.Errorf("write file %s: %w", path, err)
		}
	}
	return nil
}

type Section struct {
	Title    string
	Filename string
	Content  string
}

func extractSections(doc ast.Node, source []byte) []Section {
	var sections []Section
	var headings []*ast.Heading

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if heading, ok := n.(*ast.Heading); ok {
				if heading.Level == 2 {
					headings = append(headings, heading)
				}
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for i, heading := range headings {
		title := extractHeadingText(heading, source)
		startPos := heading.Lines().At(heading.Lines().Len() - 1).Stop
		var endPos int
		if i+1 == len(headings) {
			endPos = len(source)
		} else {
			endPos = headings[i+1].Lines().At(0).Start - 1
		}

		buf := strings.Builder{}
		buf.WriteString("# " + title + "\n\n")
		if startPos < endPos && startPos < len(source) {
			buf.Write(source[startPos:endPos])
		}
		content := dirReplacer.Replace(buf.String())
		sections = append(sections, Section{
			Title:    title,
			Filename: slugify(title) + ".md",
			Content:  strings.TrimSpace(content) + "\n",
		})
	}

	return sections
}

func extractHeadingText(heading *ast.Heading, source []byte) string {
	var text strings.Builder
	err := ast.Walk(heading, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if textNode, ok := n.(*ast.Text); ok {
				segment := textNode.Segment
				text.Write(segment.Value(source))
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return text.String()
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func checkErr(f func() error) {
	if err := f(); err != nil {
		log.Fatal(err)
	}
}
