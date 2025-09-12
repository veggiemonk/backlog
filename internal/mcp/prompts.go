package mcp

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
)

// MCPToolCall represents a structured MCP tool call that can be marshaled to JSON
type MCPToolCall struct {
	Name      string `json:"name"`
	Arguments any    `json:"arguments"`
}

// TaskCreateArgs represents arguments for task_create tool
type TaskCreateArgs struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

// TaskListArgs represents arguments for task_list tool
type TaskListArgs struct {
	Status     string   `json:"status,omitempty"`
	Sort       []string `json:"sort,omitempty"`
	Reverse    bool     `json:"reverse,omitempty"`
	Unassigned bool     `json:"unassigned,omitempty"`
	Parent     string   `json:"parent,omitempty"`
	Labels     string   `json:"labels,omitempty"`
}

// formatToolCall converts an MCPToolCall to a formatted JSON string
func formatToolCall(call MCPToolCall) string {
	bytes, _ := json.MarshalIndent(call, "", "  ")
	return string(bytes)
}

// formatMultipleToolCalls formats multiple tool calls with descriptions
func formatMultipleToolCalls(description string, calls ...MCPToolCall) string {
	result := description + "\n\n"
	for i, call := range calls {
		if i > 0 {
			result += "\n\n"
		}
		result += formatToolCall(call)
	}
	return result
}

// addPrompts adds all MCP prompts to the server
func (s *Server) addPrompts() {
	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "create_bug_report",
		Description: "Create a new bug report task.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_create",
			Arguments: core.CreateTaskParams{
				Title:       "Bug: [BUG_TITLE]",
				Description: "[BUG_DESCRIPTION]",
				Labels:      []string{"bug"},
			},
		}
		text := "Create a new bug report task using:\n" + formatToolCall(call)
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{{Content: &mcp.TextContent{Text: text}}},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "weekly_summary",
		Description: "Generate a summary of tasks completed in the last week.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Status:  []string{"done"},
				Sort:    []string{"updated"},
				Reverse: true,
			},
		}
		text := "Generate summary of completed tasks using:\n" + formatToolCall(call)
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{{Content: &mcp.TextContent{Text: text}}},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "prioritize_todo",
		Description: "Show high priority tasks that need to be done.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Status:  []string{"todo"},
				Sort:    []string{"priority"},
				Reverse: true,
			},
		}
		text := "Show high priority tasks using:\n" + formatToolCall(call)
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{{Content: &mcp.TextContent{Text: text}}},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "create_feature_request",
		Description: "Create a new feature request task.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_create",
			Arguments: core.CreateTaskParams{
				Title:       "Feature: [FEATURE_TITLE]",
				Description: "[FEATURE_DESCRIPTION]",
				Labels:      []string{"enhancement", "feature"},
			},
		}
		text := "Create a feature request using:\n" + formatToolCall(call)
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{{Content: &mcp.TextContent{Text: text}}},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "blocked_tasks",
		Description: "Find tasks that are blocked or waiting on dependencies.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		tasksDependees := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Status:        []string{"todo", "in-progress"},
				HasDependency: true,
			},
		}
		tasksDependents := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Status:     []string{"todo", "in-progress"},
				DependedOn: true,
			},
		}
		text := formatMultipleToolCalls("Find blocked tasks and check their dependencies:", tasksDependees, tasksDependents)
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{{Content: &mcp.TextContent{Text: text}}},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "project_overview",
		Description: "Get an overview of all tasks grouped by project or epic.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Parent: ptr("[PARENT_ID]"),
				Sort:   []string{"id"},
			},
		}
		text := "Get project overview using:\n" + formatToolCall(call)
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{{Content: &mcp.TextContent{Text: text}}},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "quick_standup",
		Description: "Generate a quick standup report showing what was done, what's in progress, and what's next.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		doneCall := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Status:  []string{"done"},
				Sort:    []string{"updated"},
				Reverse: true,
			},
		}
		inProgressCall := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Status: []string{"in-progress"},
			},
		}
		todoCall := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Status:  []string{"todo"},
				Sort:    []string{"priority"},
				Reverse: true,
			},
		}

		text := formatMultipleToolCalls("Generate standup report:",
			doneCall,       // What was done
			inProgressCall, // What's in progress
			todoCall)       // What's next

		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Content: &mcp.TextContent{Text: text},
				},
			},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "review_completed",
		Description: "Review recently completed tasks for retrospective or reporting.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Status:  []string{"done"},
				Sort:    []string{"updated"},
				Reverse: true,
			},
		}
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Content: &mcp.TextContent{Text: "Review completed tasks using:\n" + formatToolCall(call)},
				},
			},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "search_by_label",
		Description: "Search for tasks by a specific label or tag.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Labels: []string{"[LABEL_NAME]"},
			},
		}
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Content: &mcp.TextContent{Text: "Search tasks by label using:\n" + formatToolCall(call)},
				},
			},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "unassigned_tasks",
		Description: "Find tasks that haven't been assigned to anyone.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Unassigned: true,
				Status:     []string{"todo"},
			},
		}
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Content: &mcp.TextContent{Text: "Find unassigned tasks using:\n" + formatToolCall(call)},
				},
			},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "create_epic",
		Description: "Create a new epic or large project task that can contain subtasks.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_create",
			Arguments: core.CreateTaskParams{
				Title:       "Epic: [EPIC_TITLE]",
				Description: "[EPIC_DESCRIPTION]",
				Labels:      []string{"epic", "project"},
			},
		}
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Content: &mcp.TextContent{Text: "Create a new epic or large project task that can contain subtasks.:\n" + formatToolCall(call)},
				},
			},
		}, nil
	})

	s.mcpServer.AddPrompt(&mcp.Prompt{
		Name:        "user_story_breakdown",
		Description: "Break down a user story into smaller tasks.",
	}, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		call := MCPToolCall{
			Name: "task_list",
			Arguments: core.ListTasksParams{
				Parent: ptr("[PARENT_ID]"),
				Sort:   []string{"id"},
			},
		}
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Content: &mcp.TextContent{Text: `
Break down user story in as many tasks as needed for the whole story to be complete.
Your task is to convert the following plan into a list of tasks that can be added to a backlog using the backlog MCP tool.
Each task should be concise and actionable.

GOAL: break down the following plan into as many tasks as needed to keep the plan accurate. 

Instructions:
1. Read the plan carefully and identify all the key tasks that need to be completed.
2. For each task, provide a clear and concise title.
3. Write a detailed description of what needs to be done for each task.
4. Include Acceptance Criteria to define when the task is considered complete.
5. Outline an Implementation Plan for how the task will be carried out.
6. Add any relevant Notes that may help in completing the task.
7. Set relevant labels for each task to categorize them.
8. Assign each task to @me (you can change this later).
9. Set the status of each task to "todo".
10. Assign a priority level to each task (low, medium, high).
` + formatToolCall(call)},
				},
			},
		}, nil
	})
}

func ptr[T any](v T) *T {
	return &v
}
