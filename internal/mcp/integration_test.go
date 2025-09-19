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
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
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

func TestMCP_Integration_EndToEnd_HTTP(t *testing.T) {
	// Isolated in-memory store pre-seeded with tasks for search/list behavior
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, ".backlog")
	setupTestData(t, store)

	// Start server (autocommit disabled to avoid git dependency in tests)
	srv, err := NewServer(store, false)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	endpoint, shutdown := startHTTPServer(t, srv)
	defer func() { _ = shutdown(context.Background()) }()

	// Connect client session
	sess, closeSess := newClient(t, endpoint)
	defer func() { _ = closeSess(context.Background()) }()

	is := is.NewRelaxed(t)

	// Discover server capabilities: tools, prompts, resources
	{
		res, err := sess.ListTools(t.Context(), &mcp.ListToolsParams{})
		is.NoErr(err)
		is.Equal(len(res.Tools), 7) // task_create, task_batch_create, task_list, task_view, task_edit, task_search, task_archive
	}
	{
		res, err := sess.ListPrompts(t.Context(), &mcp.ListPromptsParams{})
		is.NoErr(err)
		is.Equal(len(res.Prompts), 12)
	}
	{
		res, err := sess.ListResources(t.Context(), &mcp.ListResourcesParams{})
		is.NoErr(err)
		is.Equal(len(res.Resources), 3)
	}

	// 1) List tasks and pick one to work with
	var picked core.Task
	{
		res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{Name: "task_list", Arguments: core.ListTasksParams{}})
		is.NoErr(err)
		wrapped := struct{ Tasks []*core.Task }{}
		parseTextContent(t, res, &wrapped)
		is.Equal(len(wrapped.Tasks), 8)
		picked = *wrapped.Tasks[0]
	}

	// 2) View the picked task
	{
		res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{Name: "task_view", Arguments: ViewParams{ID: picked.ID.String()}})
		is.NoErr(err)
		var viewed core.Task
		parseTextContent(t, res, &viewed)
		is.Equal(viewed.ID.String(), picked.ID.String())
	}

	// 3) Search for a known keyword
	{
		res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{Name: "task_search", Arguments: SearchParams{Query: "feature"}})
		is.NoErr(err)
		// Should have results thanks to setupTestData seeding "feature" labeled tasks
		wrapped := struct{ Tasks []*core.Task }{}
		parseTextContent(t, res, &wrapped)
		is.Equal(len(wrapped.Tasks), 4)
	}

	// 4) Archive the picked task
	{
		res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{Name: "task_archive", Arguments: ArchiveParams{ID: picked.ID.String()}})
		is.NoErr(err)
		is.True(res != nil)
	}

	// 7) Read a resource (AGENTS.md)
	{
		_, err := sess.ReadResource(context.Background(), &mcp.ReadResourceParams{URI: agentInstructionsURI})
		is.NoErr(err)
	}

	// 8) Get a prompt
	{
		_, err := sess.GetPrompt(context.Background(), &mcp.GetPromptParams{Name: "weekly_summary"})
		is.NoErr(err)
	}
}
