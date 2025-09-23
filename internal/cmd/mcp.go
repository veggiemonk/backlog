package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/logging"
	mcpserver "github.com/veggiemonk/backlog/internal/mcp"
)

var (
	mcpHTTPPort   int
	httpTransport bool
)

func init() {
	rootCmd.AddCommand(mcpCmd)
	setMCPFlags(mcpCmd)
}

func setMCPFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&mcpHTTPPort, "port", 8106, "Port for the MCP server (HTTP transport)")
	cmd.Flags().BoolVar(&httpTransport, "http", false, "Use HTTP transport instead of stdio")
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the MCP server",
	Long:  `Starts an MCP server to provide programmatic access to backlog tasks.`,
	Example: `
backlog mcp --http             # Start the MCP server using HTTP transport on default port 8106
backlog mcp --http --port 4321 # Start the MCP server using HTTP transport on port 4321
backlog mcp                    # Start the MCP server using stdio transport
`,
	RunE: runMCPServer,
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	store := cmd.Context().Value(ctxKeyStore).(TaskStore)
	server, err := mcpserver.NewServer(store, viper.GetBool(configAutoCommit))
	if err != nil {
		return fmt.Errorf("create MCP server: %v", err)
	}
	if httpTransport {
		logging.Info("starting MCP server", "transport", "http", "port", mcpHTTPPort)
		if err := server.RunHTTP(mcpHTTPPort); err != nil {
			return fmt.Errorf("HTTP server failed: %v", err)
		}
		return nil
	}
	if err := server.RunStdio(context.Background()); err != nil {
		return fmt.Errorf("stdio server failed: %v", err)
	}
	return nil
}
