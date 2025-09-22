package mcp

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/veggiemonk/backlog/internal/core"
)

// taskJSONSchema returns an explicit JSON schema for core.Task matching its JSON encoding.
func taskJSONSchema() *jsonschema.Schema {
	mkStringOrStringArray := func() *jsonschema.Schema {
		return &jsonschema.Schema{OneOf: []*jsonschema.Schema{
			{Type: "string"},
			{Type: "array", Items: &jsonschema.Schema{Type: "string"}},
		}}
	}
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"id":     {Type: "string"},
			"title":  {Type: "string"},
			"status": {Type: "string", Enum: []any{string(core.StatusTodo), string(core.StatusInProgress), string(core.StatusDone), string(core.StatusCancelled), string(core.StatusArchived), string(core.StatusRejected)}},
			"parent": {Type: "string"},
			"assigned": {OneOf: []*jsonschema.Schema{
				mkStringOrStringArray(), {Type: "null"},
			}},
			"labels": {OneOf: []*jsonschema.Schema{
				mkStringOrStringArray(), {Type: "null"},
			}},
			"dependencies": {OneOf: []*jsonschema.Schema{
				mkStringOrStringArray(), {Type: "null"},
			}},
			"priority":   {Type: "string", Enum: []any{"unknown", "low", "medium", "high", "critical"}},
			"created_at": {Type: "string"},
			"updated_at": {OneOf: []*jsonschema.Schema{{Type: "string"}, {Type: "null"}}},
			"history": {OneOf: []*jsonschema.Schema{
				{Type: "array", Items: &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{"timestamp": {Type: "string"}, "change": {Type: "string"}}}},
				{Type: "null"},
			}},
			"description": {Type: "string"},
			"acceptance_criteria": {OneOf: []*jsonschema.Schema{
				{Type: "array", Items: &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{"text": {Type: "string"}, "checked": {Type: "boolean"}, "index": {Type: "integer"}}}},
				{Type: "null"},
			}},
			"implementation_plan":  {OneOf: []*jsonschema.Schema{{Type: "string"}, {Type: "null"}}},
			"implementation_notes": {OneOf: []*jsonschema.Schema{{Type: "string"}, {Type: "null"}}},
		},
	}
}

// wrappedTaskJSONSchema returns a JSON schema for the wrapped Task structure
// that matches what's returned in StructuredContent: struct{ Task *core.Task }
func wrappedTaskJSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"Task": taskJSONSchema(),
		},
		Required:             []string{"Task"},
		AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}, // No additional properties allowed
	}
}

// wrappedTasksJSONSchema returns a JSON schema for the wrapped Tasks array structure
// that matches what's returned in StructuredContent: struct{ Tasks []*core.Task }
func wrappedTasksJSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"Tasks": {Type: "array", Items: taskJSONSchema()},
		},
		Required:             []string{"Tasks"},
		AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}, // No additional properties allowed
	}
}

// paginationInfoJSONSchema returns a JSON schema for core.PaginationInfo
func paginationInfoJSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"total_results":     {Type: "integer", Description: "Total number of results matching the query"},
			"displayed_results": {Type: "integer", Description: "Number of results in this response"},
			"offset":            {Type: "integer", Description: "Number of results skipped from the beginning"},
			"limit":             {Type: "integer", Description: "Maximum number of results requested (0 means no limit)"},
			"has_more":          {Type: "boolean", Description: "Whether there are more results available"},
		},
		Required:             []string{"total_results", "displayed_results", "offset", "limit", "has_more"},
		AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
	}
}

// listResultJSONSchema returns a JSON schema for core.ListResult with pagination
// that matches what's returned in StructuredContent: core.ListResult
func listResultJSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"tasks": {Type: "array", Items: taskJSONSchema()},
			"pagination": {OneOf: []*jsonschema.Schema{
				paginationInfoJSONSchema(),
				{Type: "null"},
			}},
		},
		Required:             []string{"tasks"},
		AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
	}
}
