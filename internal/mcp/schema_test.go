package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matryer/is"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/afero"
	"github.com/veggiemonk/backlog/internal/core"
)

// TestOutputSchemaCompliance verifies that the StructuredContent returned by MCP tools
// matches their declared OutputSchema.
func TestOutputSchemaCompliance(t *testing.T) {
	// Setup test store
	fs := afero.NewMemMapFs()
	store := core.NewFileTaskStore(fs, ".backlog")

	// Create a test task
	task, err := store.Create(core.CreateTaskParams{
		Title:       "Test Schema Compliance",
		Description: "Test task for schema compliance",
		Priority:    "high",
	})
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	is := is.New(t)

	// Create server
	server, err := NewServer(store, false)
	is.NoErr(err)

	t.Run("task_view schema compliance", func(t *testing.T) {
		// Call the tool
		result, _, err := server.handler.view(context.Background(), &mcp.CallToolRequest{}, ViewParams{
			ID: task.ID.String(),
		})
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)

		// Validate the StructuredContent against the expected schema
		expectedSchema := wrappedTaskJSONSchema()
		validateStructuredContent(t, expectedSchema, result.StructuredContent)
	})

	t.Run("task_create schema compliance", func(t *testing.T) {
		// Call the tool
		result, _, err := server.handler.create(context.Background(), &mcp.CallToolRequest{}, core.CreateTaskParams{
			Title:       "New Test Task",
			Description: "Another test task",
			Priority:    "medium",
		})
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)

		// Validate the StructuredContent against the expected schema
		expectedSchema := wrappedTaskJSONSchema()
		validateStructuredContent(t, expectedSchema, result.StructuredContent)
	})

	t.Run("task_edit schema compliance", func(t *testing.T) {
		// Call the tool
		result, _, err := server.handler.edit(context.Background(), &mcp.CallToolRequest{}, core.EditTaskParams{
			ID:       task.ID.String(),
			NewTitle: ptr("Updated Test Task"),
		})
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)

		// Validate the StructuredContent against the expected schema
		expectedSchema := wrappedTaskJSONSchema()
		validateStructuredContent(t, expectedSchema, result.StructuredContent)
	})

	// Test that validation actually catches schema violations
	t.Run("schema validation catches violations", func(t *testing.T) {
		// Create invalid data that doesn't match the wrapped task schema
		invalidData := map[string]any{
			"WrongField": "should not be here",
			"Task": map[string]any{
				"id": "invalid",
				// Missing required fields
			},
		}

		// This should fail validation
		expectedSchema := wrappedTaskJSONSchema()

		// Marshal and unmarshal to ensure we have proper JSON-compatible data
		jsonData, err := json.Marshal(invalidData)
		is.NoErr(err)

		var genericData any
		err = json.Unmarshal(jsonData, &genericData)
		is.NoErr(err)

		// Resolve the schema and validate
		resolved, err := expectedSchema.Resolve(nil)
		is.NoErr(err)

		// This should return an error (validation failure)
		err = resolved.Validate(genericData)
		is.True(err != nil) // Should fail validation
	})

	// Test that Tasks schema validation catches violations
	t.Run("tasks schema validation catches violations", func(t *testing.T) {
		// Create invalid data that doesn't match the wrapped tasks schema
		invalidData := map[string]any{
			"Tasks": "not an array", // Should be an array
		}

		expectedSchema := wrappedTasksJSONSchema()

		// Marshal and unmarshal to ensure we have proper JSON-compatible data
		jsonData, err := json.Marshal(invalidData)
		is.NoErr(err)

		var genericData any
		err = json.Unmarshal(jsonData, &genericData)
		is.NoErr(err)

		// Resolve the schema and validate
		resolved, err := expectedSchema.Resolve(nil)
		is.NoErr(err)

		// This should return an error (validation failure)
		err = resolved.Validate(genericData)
		is.True(err != nil) // Should fail validation
	})

	t.Run("task_list schema compliance", func(t *testing.T) {
		// Call the tool
		result, _, err := server.handler.list(context.Background(), &mcp.CallToolRequest{}, core.ListTasksParams{})
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)

		// Validate the StructuredContent against the expected schema
		expectedSchema := listResultJSONSchema()
		validateStructuredContent(t, expectedSchema, result.StructuredContent)
	})

	t.Run("task_search schema compliance", func(t *testing.T) {
		// Call the tool
		result, _, err := server.handler.search(context.Background(), &mcp.CallToolRequest{}, SearchParams{
			Query: "test",
		})
		is.NoErr(err)
		is.True(result != nil)

		// Only validate if we have StructuredContent (search can return empty results with just Content)
		if result.StructuredContent != nil {
			expectedSchema := listResultJSONSchema()
			validateStructuredContent(t, expectedSchema, result.StructuredContent)
		}
	})

	t.Run("task_batch_create schema compliance", func(t *testing.T) {
		// Call the tool
		result, _, err := server.handler.batchCreate(context.Background(), &mcp.CallToolRequest{}, ListCreateParams{
			Tasks: []core.CreateTaskParams{
				{
					Title:       "Batch Task 1",
					Description: "First batch task",
					Priority:    "low",
				},
				{
					Title:       "Batch Task 2",
					Description: "Second batch task",
					Priority:    "medium",
				},
			},
		})
		is.NoErr(err)
		is.True(result != nil)
		is.True(result.StructuredContent != nil)

		// Validate the StructuredContent against the expected schema
		expectedSchema := wrappedTasksJSONSchema()
		validateStructuredContent(t, expectedSchema, result.StructuredContent)
	})
}

// validateStructuredContent validates that the given data conforms to the JSON schema
func validateStructuredContent(t *testing.T, schema *jsonschema.Schema, data any) {
	t.Helper()
	is := is.New(t)

	// Marshal and unmarshal to ensure we have proper JSON-compatible data
	jsonData, err := json.Marshal(data)
	is.NoErr(err)

	var genericData any
	err = json.Unmarshal(jsonData, &genericData)
	is.NoErr(err)

	// Resolve the schema and validate
	resolved, err := schema.Resolve(nil)
	is.NoErr(err)

	// Validate the data against the schema
	err = resolved.Validate(genericData)
	if err != nil {
		t.Errorf("Schema validation failed: %v\nData: %s\nSchema: %+v", err, string(jsonData), schema)
	}
}
