package mcp

import (
	"testing"

	"github.com/matryer/is"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestMCP_Integration_Read_HTTP(t *testing.T) {
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
	defer func() { _ = shutdown(t.Context()) }()

	// Connect client session
	sess, closeSess := newClient(t, endpoint)
	defer func() { _ = closeSess(t.Context()) }()

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
		is := is.New(t)
		res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "task_list", Arguments: core.ListTasksParams{}})
		is.NoErr(err)
		is.True(res != nil)
		wrapped, ok := res.StructuredContent.(struct{ Tasks []*core.Task })
		is.True(ok)

		is.Equal(len(wrapped.Tasks), 8)
		picked = *wrapped.Tasks[0]
	}

	// 2) View the picked task
	{
		res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "task_view", Arguments: ViewParams{ID: picked.ID.String()}})
		is.NoErr(err)
		is.True(res != nil)
		viewed, ok := res.StructuredContent.(*core.Task)
		is.True(ok)

		is.Equal(viewed.ID.String(), picked.ID.String())
	}

	// 3) Search for a known keyword
	{
		res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "task_search", Arguments: SearchParams{Query: "feature"}})
		is.NoErr(err)
		// Should have results thanks to setupTestData seeding "feature" labeled tasks
		is.True(res != nil)
		wrapped, ok := res.StructuredContent.(struct{ Tasks []*core.Task })
		is.True(ok)

		is.Equal(len(wrapped.Tasks), 4)
	}

	// 4) Read a resource (AGENTS.md)
	{
		_, err := sess.ReadResource(t.Context(), &mcp.ReadResourceParams{URI: agentInstructionsURI})
		is.NoErr(err)
	}

	// 5) Get a prompt
	{
		_, err := sess.GetPrompt(t.Context(), &mcp.GetPromptParams{Name: "weekly_summary"})
		is.NoErr(err)
	}
}
