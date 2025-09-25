package main

//go:generate go run ./internal/tools/docgen -out ./docs/reference -format markdown

import "github.com/veggiemonk/backlog/internal/cmd"

func main() {
	cmd.Execute()
}
