package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

const viewExamples = `
  backlog view T01           # View task T01 in markdown format
  backlog view T01 --json    # View task T01 in JSON format
  backlog view T01 -j        # View task T01 in JSON format (short flag)
`

func newViewCommand(rt *runtime) *cli.Command {
	return &cli.Command{
		Name:      "view",
		Usage:     "View a task by providing its ID",
		UsageText: "backlog view <id>",
		ArgsUsage: "<id>",
		Description: "View a task by providing its ID. You can output in markdown or JSON format.\n\nExamples:\n" +
			viewExamples,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "json", Aliases: []string{"j"}, Usage: "Print JSON output"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return cli.Exit("view requires exactly one <id> argument", 1)
			}

			store := rt.store
			if store == nil {
				return fmt.Errorf("task store not initialized")
			}

			taskID := cmd.Args().First()
			task, err := store.Get(taskID)
			if err != nil {
				return fmt.Errorf("failed to view task %q: %w", taskID, err)
			}

			if cmd.Bool("json") {
				if err := json.NewEncoder(cmd.Root().Writer).Encode(task); err != nil {
					return fmt.Errorf("failed to encode JSON for task %q: %w", taskID, err)
				}
				return nil
			}

			if _, err := fmt.Fprintf(cmd.Root().Writer, "%s\n", string(task.Bytes())); err != nil {
				return fmt.Errorf("failed to write task %q contents: %w", taskID, err)
			}
			return nil
		},
	}
}
