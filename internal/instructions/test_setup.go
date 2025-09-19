package instructions

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	mcptools "github.com/veggiemonk/backlog/internal/mcp"
)

// startHTTPServer is a wrapper around the MCP testing helper
func startHTTPServer(t *testing.T, s *mcptools.Server) (string, func(context.Context) error) {
	t.Helper()
	return mcptools.StartHTTPServerForTesting(t, s)
}

// newClient is a wrapper around the MCP testing helper
func newClient(t *testing.T, endpoint string) (*mcp.ClientSession, func(context.Context) error) {
	t.Helper()
	return mcptools.NewClientForTesting(t, endpoint)
}
