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
