package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/logging"
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
	Run:  view,
}

func view(cmd *cobra.Command, args []string) {
	store := cmd.Context().Value(ctxKeyStore).(TaskStore)
	t, err := store.Get(args[0])
	if err != nil {
		logging.Error("failed to view task", "task_id", args[0], "error", err)
		os.Exit(1)
	}

	if viewJSON {
		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(t); err != nil {
			logging.Error("failed to encode JSON", "task_id", args[0], "error", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("%s\n", string(t.Bytes()))
	}
}

func setViewFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&viewJSON, "json", "j", false, "Print JSON output")
}

func init() {
	rootCmd.AddCommand(viewCmd)
	setViewFlags(viewCmd)
}
