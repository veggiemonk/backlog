package mcp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// StartHTTPServerForTesting starts an HTTP server for the given MCP Server on a random local port
// and returns the base endpoint URL and a shutdown function.
// This is exported for use in other packages' tests.
func StartHTTPServerForTesting(t *testing.T, s *Server) (string, func(context.Context) error) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	mux := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server { return s.mcpServer }, nil)
	h := &http.Server{Handler: mux}
	go func() {
		_ = h.Serve(ln)
	}()
	endpoint := fmt.Sprintf("http://%s", ln.Addr().String())
	shutdown := func(ctx context.Context) error {
		return h.Shutdown(ctx)
	}
	return endpoint, shutdown
}

// NewClientForTesting opens a client session to the given endpoint using Streamable HTTP transport.
// This is exported for use in other packages' tests.
func NewClientForTesting(t *testing.T, endpoint string) (*mcp.ClientSession, func(context.Context) error) {
	t.Helper()
	transport := &mcp.StreamableClientTransport{Endpoint: endpoint, HTTPClient: &http.Client{Timeout: 10 * time.Second}}
	client := mcp.NewClient(&mcp.Implementation{Name: "backlog-instructions-tests", Title: "tests", Version: "test"}, &mcp.ClientOptions{})
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	sess, err := client.Connect(ctx, transport, &mcp.ClientSessionOptions{})
	cancel()
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	closeFn := func(ctx context.Context) error { _ = sess.Close(); return nil }
	return sess, closeFn
}
