package cmd

// ListExamples contains all examples for the list command
var ListExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "List All Tasks",
			Command:     "backlog list",
			Comment:     "List all tasks with all columns",
		},
		{
			Description: "Filter by Status",
			Command:     "backlog list",
			Flags: map[string]string{
				"status": "todo",
			},
			Comment: "List tasks with status \"todo\"",
		},
		{
			Description: "Filter by Multiple Statuses",
			Command:     "backlog list",
			Flags: map[string]string{
				"status": "todo,in-progress",
			},
			Comment: "List tasks with status \"todo\" or \"in-progress\"",
		},
		{
			Description: "Filter by Parent",
			Command:     "backlog list",
			Flags: map[string]string{
				"parent": "12345",
			},
			Comment: "List tasks that are sub-tasks of the task with ID \"12345\"",
		},
		{
			Description: "Filter by Assigned User",
			Command:     "backlog list",
			Flags: map[string]string{
				"assigned": "alice",
			},
			Comment: "List tasks assigned to alice",
		},
		{
			Description: "Filter Unassigned Tasks",
			Command:     "backlog list",
			Flags: map[string]string{
				"unassigned": "",
			},
			Comment: "List tasks that have no one assigned",
		},
		{
			Description: "Filter by Labels",
			Command:     "backlog list",
			Flags: map[string]string{
				"labels": "bug,feature",
			},
			Comment: "List tasks containing either \"bug\" or \"feature\" labels",
		},
		{
			Description: "Filter by Priority",
			Command:     "backlog list",
			Flags: map[string]string{
				"priority": "high",
			},
			Comment: "List all high priority tasks",
		},
		{
			Description: "Filter Tasks with Dependencies",
			Command:     "backlog list",
			Flags: map[string]string{
				"has-dependency": "",
			},
			Comment: "List tasks that have at least one dependency",
		},
		{
			Description: "Filter Blocking Tasks",
			Command:     "backlog list",
			Flags: map[string]string{
				"depended-on": "",
				"status":      "todo",
			},
			Comment: "List all the blocking tasks.",
		},
		{
			Description: "Hide Extra Fields",
			Command:     "backlog list",
			Flags: map[string]string{
				"hide-extra": "",
			},
			Comment: "Hide extra fields (labels, priority, assigned)",
		},
		{
			Description: "Sort by Priority",
			Command:     "backlog list",
			Flags: map[string]string{
				"sort": "priority",
			},
			Comment: "Sort tasks by priority",
		},
		{
			Description: "Multiple Sort Fields",
			Command:     "backlog list",
			Flags: map[string]string{
				"sort": "updated,priority",
			},
			Comment: "Sort tasks by updated date, then priority",
		},
		{
			Description: "Reverse Order",
			Command:     "backlog list",
			Flags: map[string]string{
				"reverse": "",
			},
			Comment: "Reverse the order of tasks",
		},
		{
			Description: "Markdown Output",
			Command:     "backlog list",
			Flags: map[string]string{
				"markdown": "",
			},
			Comment: "List tasks in markdown format",
		},
		{
			Description: "JSON Output",
			Command:     "backlog list",
			Flags: map[string]string{
				"json": "",
			},
			Comment: "List tasks in JSON format",
		},
		{
			Description: "Pagination - Limit",
			Command:     "backlog list",
			Flags: map[string]string{
				"limit": "10",
			},
			Comment: "List first 10 tasks",
		},
		{
			Description: "Pagination - Limit and Offset",
			Command:     "backlog list",
			Flags: map[string]string{
				"limit":  "5",
				"offset": "10",
			},
			Comment: "List 5 tasks starting from 11th task",
		},
	},
}
