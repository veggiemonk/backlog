package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/logging"
	mcpserver "github.com/veggiemonk/backlog/internal/mcp"
)

const mcpExamples = `
backlog mcp --http             # Start the MCP server using HTTP transport on default port 8106
backlog mcp --http --port 4321 # Start the MCP server using HTTP transport on port 4321
backlog mcp                    # Start the MCP server using stdio transport
`

func newMCPCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:        "mcp",
		Usage:       "Start the MCP server",
		Description: "Starts an MCP server to provide programmatic access to backlog tasks.\n\nExamples:\n" + mcpExamples,
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "port", Value: 8106, Usage: "Port for the MCP server (HTTP transport)"},
			&cli.BoolFlag{Name: "http", Usage: "Use HTTP transport instead of stdio"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			store := rt.store
			if store == nil {
				return fmt.Errorf("task store not initialized")
			}

			server, err := mcpserver.NewServer(store, rt.autoCommit)
			if err != nil {
				return fmt.Errorf("create MCP server: %v", err)
			}

			if cmd.Bool("http") {
				port := cmd.Int("port")
				logging.Info("starting MCP server", "transport", "http", "port", port)
				if err := server.RunHTTP(ctx, port); err != nil {
					return fmt.Errorf("HTTP server failed: %v", err)
				}
				return nil
			}

			if err := server.RunStdio(ctx); err != nil {
				return fmt.Errorf("stdio server failed: %v", err)
			}
			return nil
		},
	}
}
