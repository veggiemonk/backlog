package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// startHTTPServer starts an HTTP server for the given MCP Server on a random local port
// and returns the base endpoint URL and a shutdown function.
func startHTTPServer(t *testing.T, s *Server) (string, func(context.Context) error) {
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

// newClient opens a client session to the given endpoint using Streamable HTTP transport.
func newClient(t *testing.T, endpoint string) (*mcp.ClientSession, func(context.Context) error) {
	t.Helper()
	transport := &mcp.StreamableClientTransport{Endpoint: endpoint, HTTPClient: &http.Client{Timeout: 10 * time.Second}}
	client := mcp.NewClient(&mcp.Implementation{Name: "backlog-integration-tests", Title: "tests", Version: "test"}, &mcp.ClientOptions{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	sess, err := client.Connect(ctx, transport, &mcp.ClientSessionOptions{})
	cancel()
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	closeFn := func(ctx context.Context) error { _ = sess.Close(); return nil }
	return sess, closeFn
}

// parseTextContent decodes the first text content payload from a CallToolResult into the provided target.
func parseTextContent(t *testing.T, res *mcp.CallToolResult, target any) {
	t.Helper()
	is := is.New(t)
	is.True(res != nil)
	is.True(len(res.Content) > 0)
	txt, ok := res.Content[0].(*mcp.TextContent)
	is.True(ok)
	is.NoErr(json.Unmarshal([]byte(txt.Text), target))
}
