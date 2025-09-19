package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

func TestMCP_Integration_Create_Edit_Batch_HTTP(t *testing.T) {
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
	defer func() { _ = shutdown(context.Background()) }()

	// Connect client session
	sess, closeSess := newClient(t, endpoint)
	defer func() { _ = closeSess(context.Background()) }()

	is := is.New(t)

	// 1) Create a task (priority is a string enum)
	var created core.Task
	{
		params := core.CreateTaskParams{Title: "Schema Task", Description: "created via MCP", Labels: []string{"schema"}, Priority: "high"}
		res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{Name: "task_create", Arguments: params})
		is.NoErr(err)
		// Ensure server populated text content
		parseTextContent(t, res, &created)
		is.Equal(created.Title, "Schema Task")
		is.Equal(created.Priority.String(), "high")
		// Structured content should also be present (validated by client)
		is.True(res.StructuredContent != nil)
	}

	// 2) Edit the task title and verify via view
	{
		newTitle := "Schema Task (edited)"
		res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{Name: "task_edit", Arguments: core.EditTaskParams{ID: created.ID.String(), NewTitle: &newTitle}})
		is.NoErr(err)
		// Text content should be auto-populated from structured output
		var updated core.Task
		parseTextContent(t, res, &updated)
		is.Equal(updated.Title, newTitle)
		created = updated
	}

	// 3) Batch create three tasks and assert count and types
	{
		params := ListCreateParams{Tasks: []core.CreateTaskParams{
			{Title: "Batch One", Priority: "low"},
			{Title: "Batch Two", Priority: "medium"},
			{Title: "Batch Three", Priority: "critical"},
		}}
		res, err := sess.CallTool(context.Background(), &mcp.CallToolParams{Name: "task_batch_create", Arguments: params})
		is.NoErr(err)
		// Validate text content wrapper
		wrapped := struct{ Tasks []*core.Task }{}
		parseTextContent(t, res, &wrapped)
		is.Equal(len(wrapped.Tasks), 3)
		// Validate structured content exists as well
		is.True(res.StructuredContent != nil)
		// Parse structuredContent back to the same wrapper to ensure shape
		b, err := json.Marshal(res.StructuredContent)
		is.NoErr(err)
		wrapped2 := struct{ Tasks []*core.Task }{}
		is.NoErr(json.Unmarshal(b, &wrapped2))
		is.Equal(len(wrapped2.Tasks), 3)
	}
}
