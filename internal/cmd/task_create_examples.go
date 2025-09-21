package cmd

// CreateExamples contains all examples for the create command
var CreateExamples = CommandExamples{
	Examples: []CommandExample{
		{
			Description: "Basic Task Creation",
			Command:     "backlog create",
			Args:        []string{"Fix the login button styling"},
			Comment:     "This is the simplest way to create a task, providing only a title.",
		},
		{
			Description: "Task with Description",
			Command:     "backlog create",
			Args:        []string{"Implement password reset"},
			Flags: map[string]string{
				"description": "Users should be able to request a password reset link via their email. This involves creating a new API endpoint and a front-end form.",
			},
			Comment: "Use the -d or --description flag to add more detailed information about the task.",
		},
		{
			Description: "Assigning to Single Person",
			Command:     "backlog create",
			Args:        []string{"Design the new dashboard"},
			Flags: map[string]string{
				"assigned": "alex",
			},
			Comment: "You can assign a task to one or more team members using the -a or --assigned flag.",
		},
		{
			Description: "Assigning to Multiple People",
			Command:     "backlog create",
			Args:        []string{"Code review for the payment gateway"},
			Flags: map[string]string{
				"assigned": "jordan,casey",
			},
		},
		{
			Description: "Adding Labels",
			Command:     "backlog create",
			Args:        []string{"Update third-party dependencies"},
			Flags: map[string]string{
				"labels": "bug,backend,security",
			},
			Comment: "Use the -l or --labels flag to categorize the task with comma-separated labels.",
		},
		{
			Description: "Setting High Priority",
			Command:     "backlog create",
			Args:        []string{"Hotfix: Production database is down"},
			Flags: map[string]string{
				"priority": "high",
			},
			Comment: "Specify the task's priority with the --priority flag. The default is \"medium\".",
		},
		{
			Description: "Setting Low Priority",
			Command:     "backlog create",
			Args:        []string{"Refactor the old user model"},
			Flags: map[string]string{
				"priority": "low",
			},
		},
		{
			Description: "Defining Acceptance Criteria",
			Command:     "backlog create",
			Args:        []string{"Develop user profile page"},
			Flags: map[string]string{
				"ac": "Users can view their own profile information.,Users can upload a new profile picture.,The page is responsive on mobile devices.",
			},
			Comment: "Use the --ac flag multiple times to list the conditions that must be met for the task to be considered complete.",
		},
		{
			Description: "Creating a Sub-task",
			Command:     "backlog create",
			Args:        []string{"Add Google OAuth login"},
			Flags: map[string]string{
				"parent": "15",
			},
			Comment: "Link a new task to a parent task using the -p or --parent flag. This is useful for breaking down larger tasks.",
		},
		{
			Description: "Setting Single Dependency",
			Command:     "backlog create",
			Args:        []string{"Deploy user authentication"},
			Flags: map[string]string{
				"deps": "T15",
			},
			Comment: "Use the --deps flag to specify that this task depends on other tasks being completed first.",
		},
		{
			Description: "Setting Multiple Dependencies",
			Command:     "backlog create",
			Args:        []string{"Integration testing"},
			Flags: map[string]string{
				"deps": "T15,T18,T20",
			},
			Comment: "This means the task cannot be started until tasks T15, T18, and T20 are completed.",
		},
		{
			Description: "Task with Implementation Notes",
			Command:     "backlog create",
			Args:        []string{"Optimize database queries"},
			Flags: map[string]string{
				"notes": "Focus on the user lookup queries in the authentication module. Consider adding indexes on email and username fields.",
			},
			Comment: "Use the --notes flag to add implementation notes to help with development.",
		},
		{
			Description: "Task with Implementation Plan",
			Command:     "backlog create",
			Args:        []string{"Implement user registration flow"},
			Flags: map[string]string{
				"plan": "1. Design registration form UI\n2. Create user validation logic\n3. Set up email verification\n4. Add password strength requirements\n5. Write integration tests",
			},
			Comment: "Use the --plan flag to add a structured implementation plan.",
		},
		{
			Description: "Complex Example with Multiple Flags",
			Command:     "backlog create",
			Args:        []string{"Build the new reporting feature"},
			Flags: map[string]string{
				"description": "Create a new section in the app that allows users to generate and export monthly performance reports in PDF format.",
				"assigned":    "drew",
				"labels":      "feature,frontend,backend",
				"priority":    "high",
				"ac":          "Report generation logic is accurate.,Users can select a date range for the report.,The exported PDF has the correct branding and layout.",
				"parent":      "23",
			},
			Comment: "Here is a comprehensive example that uses several flags at once to create a very detailed task.",
		},
	},
}
