package cmd

import (
	"context"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

var (
	tasksDir   string
	autoCommit bool
)

type contextKey string

const (
	ctxKeyStore = contextKey("store")
	defaultDir  = ".backlog"
)

type TaskStore interface {
	Get(id string) (*core.Task, error)
	Create(params core.CreateTaskParams) (*core.Task, error)
	Update(task *core.Task, params core.EditTaskParams) (*core.Task, error)
	List(params core.ListTasksParams) ([]*core.Task, error)
	Search(query string, listParams core.ListTasksParams) ([]*core.Task, error)
	Path(t *core.Task) string
	Archive(id core.TaskID) (string, error)
}

var _ TaskStore = (*core.FileTaskStore)(nil)

var rootCmd = &cobra.Command{
	Use:   "backlog",
	Short: "Backlog is a git-native, markdown-based task manager",
	Long: `A Git-native, Markdown-based task manager for developers and AI agents.
Backlog helps you manage tasks within your git repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is provided
		if err := cmd.Help(); err != nil {
			logging.Error("failed to display help", "error", err)
			os.Exit(1)
		}
	},
}

func setRootPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&tasksDir, "folder", defaultDir, "Directory for backlog tasks")
	cmd.PersistentFlags().BoolVar(&autoCommit, "auto-commit", true, "Auto-committing changes to git repository")
}

func init() {
	setRootPersistentFlags(rootCmd)
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Initialize logging before anything else
		logging.Init()

		fs := afero.NewOsFs()
		var store TaskStore = core.NewFileTaskStore(fs, tasksDir)
		cmd.SetContext(context.WithValue(cmd.Context(), ctxKeyStore, store))
	}
}

func Root() *cobra.Command {
	return rootCmd
}

func Execute() {
	defer func() { logging.Close() }()
	if err := rootCmd.Execute(); err != nil {
		logging.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
