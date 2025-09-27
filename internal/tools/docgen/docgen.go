// Package main provides documentation generation utilities for the backlog CLI.
// It generates markdown, man pages, or restructured text documentation from
// cobra commands.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
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
	format := flag.String("format", "markdown", "markdown|man|rest")
	flag.Parse()

	checkErr(func() error { return os.MkdirAll(dirDocs, 0o750) })
	checkErr(func() error { return genReference(dirRef, *format) })
	checkErr(splitReadMe)
	checkErr(addPromptCLI)
	checkErr(addPromptMCP)
	checkErr(addAGENTSmd)
	checkErr(addIndex)
	checkErr(addEnvVars)
}

func addEnvVars() error {
	files, err := os.ReadDir(dirRef)
	if err != nil {
		return err
	}
	sep := "### Options"
	var buf strings.Builder
	buf.WriteString("\n#### Environment Variables\n\n")
	buf.WriteString("```\n")
	buf.WriteString("\t(name)\t\t(default)\n")
	m := viper.AllSettings()
	for _, k := range slices.Sorted(maps.Keys(m)) {
		if len(k) < 8 {
			buf.WriteString(fmt.Sprintf("\t%s\t\t%v\n", strings.ToUpper(k), m[k]))
		} else {
			buf.WriteString(fmt.Sprintf("\t%s\t%v\n", strings.ToUpper(k), m[k]))
		}
	}
	buf.WriteString("```\n")
	buf.WriteString("\n#### Flags\n")
	for _, f := range files {
		if filepath.Ext(f.Name()) != ".md" {
			continue
		}
		path := filepath.Join(dirRef, f.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		before, after, found := strings.Cut(string(b), sep)
		if !found {
			continue
		}
		results := fmt.Sprintf("%s\n%s\n%s%s", before, sep, buf.String(), after)

		if err := os.WriteFile(path, []byte(results), os.ModePerm); err != nil {
			return err
		}

	}

	return nil
}

func checkErr(f func() error) {
	if err := f(); err != nil {
		log.Fatal(err)
	}
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
	// 1. Parse README markdown
	b, err := os.ReadFile("README.md")
	if err != nil {
		return err
	}
	p := goldmark.DefaultParser()
	doc := p.Parse(text.NewReader(b))
	// 2. Split sections by headings
	sections := extractSections(doc, b)
	// 3. Write each section to ./docs/name_of_section.md
	for _, section := range sections {
		if section.Title == "" {
			continue // Skip sections without titles
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

	// First pass: collect all headings
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

	// Second pass: extract content between headings
	for i, heading := range headings {
		title := extractHeadingText(heading, source)

		// Find the start and end positions for this section
		startPos := heading.Lines().At(heading.Lines().Len() - 1).Stop
		var endPos int
		if i+1 == len(headings) { // last heading
			endPos = len(source)
		} else {
			endPos = headings[i+1].Lines().At(0).Start - 1
		}

		// Extract the section content
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
	// Convert to lowercase and replace spaces and special characters with underscores
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	// Remove or replace other special characters
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func genReference(out, format string) error {
	if err := os.MkdirAll(out, 0o750); err != nil {
		return err
	}

	root := cmd.Root()
	root.DisableAutoGenTag = true // stable, reproducible files (no timestamp footer)

	switch format {
	case "markdown":
		prep := func(filename string) string {
			base := filepath.Base(filename)
			name := strings.TrimSuffix(base, filepath.Ext(base))
			title := strings.ReplaceAll(name, "_", " ")
			var buf bytes.Buffer
			buf.WriteString("---\n")
			buf.WriteString("layout: page\n")
			buf.WriteString("title: " + title + "\n")
			buf.WriteString("---\n\n")
			return buf.String()
		}
		link := func(name string) string { return strings.ToLower(name) }
		if err := doc.GenMarkdownTreeCustom(root, out, prep, link); err != nil {
			return err
		}
		return nil
	case "man":
		hdr := &doc.GenManHeader{Title: strings.ToUpper(root.Name()), Section: "1"}
		if err := doc.GenManTree(root, hdr, out); err != nil {
			return err
		}
	case "rest":
		if err := doc.GenReSTTree(root, out); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
	return nil
}
