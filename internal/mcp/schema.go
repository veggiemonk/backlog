package mcp

import (
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/veggiemonk/backlog/internal/core"
)

// taskJSONSchema returns an explicit JSON schema for core.Task matching its JSON encoding.
func taskJSONSchema() *jsonschema.Schema {
	schema, _ := jsonschema.For[core.Task](&jsonschema.ForOptions{
		TypeSchemas: map[reflect.Type]*jsonschema.Schema{
			reflect.TypeFor[core.TaskID](): {OneOf: []*jsonschema.Schema{
				{Type: "string"}, {Type: "null"},
			}},
			reflect.TypeFor[core.MaybeStringArray](): {OneOf: []*jsonschema.Schema{
				{Type: "null"},
				{Type: "string"},
				{Type: "array", Items: &jsonschema.Schema{Type: "string"}},
			}},
			reflect.TypeFor[core.Priority](): {Type: "string"},
		},
	})
	return schema
}

// wrappedTasksJSONSchema returns a JSON schema for the wrapped Tasks array structure
// that matches what's returned in StructuredContent: struct{ Tasks []core.Task }
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
	schema, _ := jsonschema.For[core.PaginationInfo](nil)
	return schema
}

// listResultJSONSchema returns a JSON schema for core.ListResult with pagination
// that matches what's returned in StructuredContent: core.ListResult
func listResultJSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"tasks": {Type: "array", Items: taskJSONSchema()},
			"pagination": {
				OneOf: []*jsonschema.Schema{
					paginationInfoJSONSchema(),
					{Type: "null"},
				},
			},
		},
		Required:             []string{"tasks"},
		AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
	}
}
