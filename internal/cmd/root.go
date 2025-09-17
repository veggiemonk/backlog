// Package cmd contains the cobra commands
package cmd

import (
	"context"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
	"github.com/veggiemonk/backlog/internal/paths"
)

var (
	tasksDir   string
	autoCommit bool
)

type contextKey string

const ctxKeyStore = contextKey("store")

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

func init() {
	cobra.OnInitialize(initConfig)
	setRootPersistentFlags(rootCmd)
	rootCmd.PersistentPreRun = preRun
}

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

func initConfig() {
	// Set environment variable prefix
	viper.SetEnvPrefix("BACKLOG")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("folder", paths.DefaultDir)
	viper.SetDefault("auto-commit", true)
	viper.SetDefault("log-level", "info")
	viper.SetDefault("log-format", "text")
	viper.SetDefault("log-file", "")

	// Bind environment variables with their keys
	viper.BindEnv("folder", "BACKLOG_FOLDER")
	viper.BindEnv("auto-commit", "BACKLOG_AUTO_COMMIT")
	viper.BindEnv("log-level", "BACKLOG_LOG_LEVEL")
	viper.BindEnv("log-format", "BACKLOG_LOG_FORMAT")
	viper.BindEnv("log-file", "BACKLOG_LOG_FILE")
}

func setRootPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&tasksDir, "folder", paths.DefaultDir, "Directory for backlog tasks")
	cmd.PersistentFlags().BoolVar(&autoCommit, "auto-commit", true, "Auto-committing changes to git repository")
	cmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	cmd.PersistentFlags().String("log-format", "text", "Log format (json, text)")
	cmd.PersistentFlags().String("log-file", "", "Log file path (defaults to stderr)")

	// Bind flags to viper
	viper.BindPFlag("folder", cmd.PersistentFlags().Lookup("folder"))
	viper.BindPFlag("auto-commit", cmd.PersistentFlags().Lookup("auto-commit"))
	viper.BindPFlag("log-level", cmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("log-format", cmd.PersistentFlags().Lookup("log-format"))
	viper.BindPFlag("log-file", cmd.PersistentFlags().Lookup("log-file"))
}

func preRun(cmd *cobra.Command, args []string) {
	// Initialize logging using Viper values
	logging.Init(
		viper.GetString("log-level"),
		viper.GetString("log-format"),
		viper.GetString("log-file"),
	)

	// Use Viper to get the tasks directory
	tasksDir = viper.GetString("folder")
	autoCommit = viper.GetBool("auto-commit")

	logging.Debug("resolve env var", "tasksDir", tasksDir, "autoCommit", autoCommit)
	fs := afero.NewOsFs()
	var err error
	tasksDir, err = paths.ResolveTasksDir(fs, tasksDir)
	if err != nil {
		logging.Error("tasks directory", "error", err)
	}
	logging.Debug("resolve tasks directory", "tasksDir", tasksDir)
	var store TaskStore = core.NewFileTaskStore(fs, tasksDir)
	cmd.SetContext(context.WithValue(cmd.Context(), ctxKeyStore, store))
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
