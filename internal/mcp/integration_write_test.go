package mcp

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestMCP_Integration_Write_HTTP(t *testing.T) {
	t.Parallel()
	// Setup isolated in-memory store and seed
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, ".backlog")
	setupTestData(t, store)

	// Start server (with explicit OutputSchemas from code)
	srv, err := NewServer(store, false)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	endpoint, shutdown := startHTTPServer(t, srv)
	defer func() { _ = shutdown(t.Context()) }()

	// Connect client session
	sess, closeSess := newClient(t, endpoint)
	defer func() { _ = closeSess(t.Context()) }()

	is := is.New(t)

	// 1) Create a task (priority is a string enum)
	var created *core.Task
	{
		params := core.CreateTaskParams{Title: "Schema Task", Description: "created via MCP", Labels: []string{"schema"}, Priority: "high"}
		res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "task_create", Arguments: params})
		is.NoErr(err)
		is.True(res != nil)
		// Structured content comes back as generic JSON over HTTP; decode it
		wrapped := core.Task{}
		b, err := json.Marshal(res.StructuredContent)
		is.NoErr(err)
		is.NoErr(json.Unmarshal(b, &wrapped))
		created = &wrapped

		is.Equal(created.Title, "Schema Task")
		is.Equal(created.Priority.String(), "high")
	}

	// 2) Edit the task title and verify via view
	{
		newTitle := "Schema Task (edited)"
		params := core.EditTaskParams{ID: created.ID.String(), NewTitle: &newTitle}
		res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "task_edit", Arguments: params})
		is.NoErr(err)
		// Decode structured content
		wrapped := core.Task{}
		b, err := json.Marshal(res.StructuredContent)
		is.NoErr(err)
		is.NoErr(json.Unmarshal(b, &wrapped))

		is.Equal(wrapped.Title, newTitle)
		created = &wrapped
	}

	// 3) Batch create three tasks and assert count and types
	{
		params := ListCreateParams{Tasks: []core.CreateTaskParams{
			{Title: "Batch One", Priority: "low"},
			{Title: "Batch Two", Priority: "medium"},
			{Title: "Batch Three", Priority: "critical"},
		}}
		res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "task_batch_create", Arguments: params})
		is.NoErr(err)
		is.True(res != nil)
		wrapped := []core.Task{}
		b, err := json.Marshal(res.StructuredContent)
		is.NoErr(err)
		is.NoErr(json.Unmarshal(b, &wrapped))

		is.Equal(len(wrapped), 3)
		// Parse structuredContent back to the same wrapper to ensure shape
		b, err = json.Marshal(res.StructuredContent)
		is.NoErr(err)
		wrapped2 := []core.Task{}
		is.NoErr(json.Unmarshal(b, &wrapped2))
		is.Equal(len(wrapped2), 3)
	}
	// 4) Archive the picked task
	{
		res, err := sess.CallTool(t.Context(), &mcp.CallToolParams{Name: "task_archive", Arguments: ArchiveParams{ID: created.ID.String()}})
		is.NoErr(err)
		is.True(res != nil)
	}
}
