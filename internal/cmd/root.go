// Package cmd contains the cobra commands
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
	mcpserver "github.com/veggiemonk/backlog/internal/mcp"
	"github.com/veggiemonk/backlog/internal/paths"
)

type contextKey string

const ctxKeyStore = contextKey("store")

var _ mcpserver.TaskStore = (*core.FileTaskStore)(nil)

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
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default action when no subcommand is provided
		if err := cmd.Help(); err != nil {
			return fmt.Errorf("failed to display help: %w", err)
		}
		return nil
	},
}

const (
	// Environment variable names for configuration
	envPrefix        = "BACKLOG"
	envVarLogFile    = "BACKLOG_LOG_FILE"
	envVarLogLevel   = "BACKLOG_LOG_LEVEL"
	envVarLogFormat  = "BACKLOG_LOG_FORMAT"
	envVarAutoCommit = "BACKLOG_AUTO_COMMIT"

	// folder
	configFolder  = "folder"
	envVarDir     = "BACKLOG_FOLDER"
	defaultFolder = ".backlog"

	// git
	configAutoCommit  = "auto-commit"
	defaultAutoCommit = true

	// logging
	configLogLevel   = "log-level"
	defaultLogLevel  = "info"
	configLogFormat  = "log-format"
	defaultLogFormat = "text"
	configLogFile    = "log-file"
	defaultLogFile   = ""
)

func preRun(cmd *cobra.Command, args []string) {
	// Initialize logging using Viper values
	logging.Init(
		viper.GetString(configLogLevel),
		viper.GetString(configLogFormat),
		viper.GetString(configLogFile),
	)

	// Use Viper to get the tasks directory
	tasksDir := viper.GetString(configFolder)
	autoCommit := viper.GetBool(configAutoCommit)

	logging.Debug("resolve env var", configFolder, tasksDir, configAutoCommit, autoCommit)
	fs := afero.NewOsFs()
	var err error
	tasksDir, err = paths.ResolveTasksDir(fs, tasksDir)
	if err != nil {
		logging.Error("tasks directory", "error", err)
	}
	logging.Debug("resolve tasks directory", configFolder, tasksDir)
	var store mcpserver.TaskStore = core.NewFileTaskStore(fs, tasksDir)
	cmd.SetContext(context.WithValue(cmd.Context(), ctxKeyStore, store))
}

func initConfig() {
	// Set environment variable prefix
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault(configFolder, defaultFolder)
	viper.SetDefault(configAutoCommit, defaultAutoCommit)
	viper.SetDefault(configLogLevel, defaultLogLevel)
	viper.SetDefault(configLogFormat, defaultLogFormat)
	viper.SetDefault(configLogFile, defaultLogFile)

	// Bind environment variables with their keys
	checkErr(viper.BindEnv(configFolder, envVarDir))
	checkErr(viper.BindEnv(configAutoCommit, envVarAutoCommit))
	checkErr(viper.BindEnv(configLogLevel, envVarLogLevel))
	checkErr(viper.BindEnv(configLogFormat, envVarLogFormat))
	checkErr(viper.BindEnv(configLogFile, envVarLogFile))
}

func checkErr(err error) {
	if err != nil {
		logging.Error("binding environment variables", "err", err)
	}
}

func setRootPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(configFolder, defaultFolder, "Directory for backlog tasks")
	cmd.PersistentFlags().Bool(configAutoCommit, defaultAutoCommit, "Auto-committing changes to git repository")
	cmd.PersistentFlags().String(configLogLevel, defaultLogLevel, "Log level (debug, info, warn, error)")
	cmd.PersistentFlags().String(configLogFormat, defaultLogFormat, "Log format (json, text)")
	cmd.PersistentFlags().String(configLogFile, defaultLogFile, "Log file path (defaults to stderr)")

	// Bind flags to viper
	checkErr(viper.BindPFlag(configFolder, cmd.PersistentFlags().Lookup(configFolder)))
	checkErr(viper.BindPFlag(configAutoCommit, cmd.PersistentFlags().Lookup(configAutoCommit)))
	checkErr(viper.BindPFlag(configLogLevel, cmd.PersistentFlags().Lookup(configLogLevel)))
	checkErr(viper.BindPFlag(configLogFormat, cmd.PersistentFlags().Lookup(configLogFormat)))
	checkErr(viper.BindPFlag(configLogFile, cmd.PersistentFlags().Lookup(configLogFile)))
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