## backlog list

List all tasks

### Synopsis

Lists all tasks in the backlog except archived tasks.

```
backlog list [flags]
```

### Examples

```

# List all tasks
backlog list 							# List all tasks with all columns
backlog list --status "todo" 			# List tasks with status "todo"
backlog list --status "todo,in-progress" # List tasks with status "todo" or "in-progress"
backlog list --status "done" 			# List tasks with status "done"
backlog list --parent "12345" 			# List tasks that are sub-tasks of the task with ID "12345"
backlog list --status "todo" --parent "12345" # List "todo" sub-tasks of task "12345"
backlog list --assigned "alice" 		# List tasks assigned to alice
backlog list --unassigned 				# List tasks that have no one assigned
backlog list --labels "bug" 			# List tasks containing the label "bug"
backlog list --labels "bug,feature" 	# List tasks containing either "bug" or "feature" labels

# dependency filters
backlog list --has-dependency 	# List tasks that have at least one dependency
backlog list --depended-on 		# List tasks that are depended on by other tasks
backlog list --depended-on --status "todo" 	# List all the blocking tasks.

# column visibility
backlog list --hide-extra 		# Hide extra fields (labels, priority, assigned)
backlog list -e 				# Hide extra fields (labels, priority, assigned)
backlog --status "todo" --hide-extra # List "todo" tasks with minimal columns

# sorting
backlog list --sort "priority" 			# Sort tasks by priority
backlog list --sort "updated,priority" 	# Sort tasks by updated date, then priority
backlog list --sort "status,created" 	# Sort tasks by status, then creation date
backlog list --reverse 					# Reverse the order of tasks
backlog list --sort "priority" --reverse 	# Sort by priority in reverse order
backlog list --status "todo" --sort "priority" --reverse # Combine all options

# output format
backlog list -m 		# List tasks in markdown format
backlog list -markdown 	# List tasks in markdown format
backlog list --json 	# List tasks in JSON format
backlog list -j  		# List tasks in JSON format
backlog list --status "todo" --json # List "todo" tasks in JSON format
	
```

### Options

```
  -a, --assigned strings   Filter tasks by assigned names
  -d, --depended-on        Filter tasks that are depended on by other tasks
  -c, --has-dependency     Filter tasks that have dependencies
  -h, --help               help for list
  -e, --hide-extra         Hide extra fields (labels, priority, assigned)
  -j, --json               Print JSON output
  -l, --labels strings     Filter tasks by labels
  -m, --markdown           print markdown table
  -p, --parent string      Filter tasks by parent ID
  -r, --reverse            Reverse the order of tasks
      --sort string        Sort tasks by comma-separated fields (id, title, status, priority, created, updated)
  -s, --status strings     Filter tasks by status
  -u, --unassigned         Filter tasks that have no one assigned
```

### Options inherited from parent commands

```
      --auto-commit     Auto-committing changes to git repository (default true)
      --folder string   Directory for backlog tasks (default ".backlog")
```

### SEE ALSO

* [backlog](backlog.md)	 - Backlog is a git-native, markdown-based task manager

