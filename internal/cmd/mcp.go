package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/logging"
	mcpserver "github.com/veggiemonk/backlog/internal/mcp"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the MCP server",
	Long:  `Starts an MCP server to provide programmatic access to backlog tasks.`,
	Example: `
backlog mcp --http             # Start the MCP server using HTTP transport on default port 8106
backlog mcp --http --port 4321 # Start the MCP server using HTTP transport on port 4321
backlog mcp                    # Start the MCP server using stdio transport
`,
	Run: runMcpServer,
}

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

func runMcpServer(cmd *cobra.Command, args []string) {
	store := cmd.Context().Value(ctxKeyStore).(TaskStore)
	server := mcpserver.NewServer(store, viper.GetBool(configAutoCommit))
	if httpTransport {
		logging.Info("starting MCP server", "transport", "http", "port", mcpHTTPPort)
		if err := server.RunHTTP(mcpHTTPPort); err != nil {
			logging.Error("HTTP server failed", "error", err)
		}
		return
	}
	if err := server.RunStdio(context.Background()); err != nil {
		logging.Error("stdio server failed", "error", err)
	}
}
