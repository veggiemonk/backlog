package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	mcpserver "github.com/veggiemonk/backlog/internal/mcp"
)

var viewJSON bool

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View a task by providing its ID",
	Long: `View a task by providing its ID. You can output in markdown or JSON format.

Examples:
  backlog view T01           # View task T01 in markdown format
  backlog view T01 --json    # View task T01 in JSON format
  backlog view T01 -j        # View task T01 in JSON format (short flag)`,
	Args: cobra.ExactArgs(1),
	RunE: view,
}

func view(cmd *cobra.Command, args []string) error {
	store := cmd.Context().Value(ctxKeyStore).(mcpserver.TaskStore)
	t, err := store.Get(args[0])
	if err != nil {
		return fmt.Errorf("failed to view task %q: %w", args[0], err)
	}

	if viewJSON {
		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(t); err != nil {
			return fmt.Errorf("failed to encode JSON for task %q: %w", args[0], err)
		}
	} else {
		fmt.Printf("%s\n", string(t.Bytes()))
	}
	return nil
}

func setViewFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&viewJSON, "json", "j", false, "Print JSON output")
}

func init() {
	rootCmd.AddCommand(viewCmd)
	setViewFlags(viewCmd)
}
