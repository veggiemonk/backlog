package instructions

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/veggiemonk/backlog/internal/core"
)

type Generic[In, Out any] struct {
	MCPName     string
	CLIName     string
	Description string
	IO          struct {
		In  In
		Out Out
	}
	Schema struct {
		In  *jsonschema.Schema
		Out *jsonschema.Schema
	}
	Examples []Example[In, Out]
}

type Schema struct {
	In  *jsonschema.Schema
	Out *jsonschema.Schema
}
type IO[In, Out any] struct {
	In  In
	Out Out
}
type Example[In, Out any] struct {
	Name        string
	Description string
	Params      In
	Expected    Out
}

var Create = Generic[core.CreateTaskParams, core.Task]{
	MCPName:     "task_create",
	CLIName:     "create",
	Description: `# Create tasks using the "backlog create" command with its different flags.`,
	Schema:      Schema{In: nil, Out: nil},
	Examples: []Example[core.CreateTaskParams, core.Task]{
		{
			Name:        "Basic Task Creation",
			Description: "This is the simplest way to create a task, providing only a title.",
			Params:      core.CreateTaskParams{Title: `"Fix the login button styling"`},
			Expected:    core.Task{Title: "Fix the login button styling"},
		},
	},
}
