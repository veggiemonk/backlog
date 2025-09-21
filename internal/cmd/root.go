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
	"github.com/veggiemonk/backlog/internal/paths"
)

const Name = "backlog"

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
	Use:   Name,
	Short: "git-native, markdown-based task manager",
	Long: `A Git-native, Markdown-based task manager for developers and AI agents.
Backlog helps you manage tasks within your git repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default action when no subcommand is provided
		if err := cmd.Help(); err != nil {
			return fmt.Errorf("display help: %v", err)
		}
		return nil
	},
}

const (
	// Environment variable names for configuration
	envPrefix        = "BACKLOG"
	envVarLogFile    = envPrefix + "_LOG_FILE"
	envVarLogLevel   = envPrefix + "_LOG_LEVEL"
	envVarLogFormat  = envPrefix + "_LOG_FORMAT"
	envVarAutoCommit = envPrefix + "_AUTO_COMMIT"
	envVarPageSize   = envPrefix + "_PAGE_SIZE"
	envVarMaxLimit   = envPrefix + "_MAX_LIMIT"

	// folder
	configFolder  = "folder"
	envVarDir     = envPrefix + "_FOLDER"
	defaultFolder = ".backlog"

	// git
	configAutoCommit  = "auto-commit"
	defaultAutoCommit = false

	// pagination
	configPageSize  = "page-size"
	defaultPageSize = 25
	configMaxLimit  = "max-limit"
	defaultMaxLimit = 1000

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
	var store TaskStore = core.NewFileTaskStore(fs, tasksDir)
	cmd.SetContext(context.WithValue(cmd.Context(), ctxKeyStore, store))
}

func initConfig() {
	// Set environment variable prefix
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault(configFolder, defaultFolder)
	viper.SetDefault(configAutoCommit, defaultAutoCommit)
	viper.SetDefault(configPageSize, defaultPageSize)
	viper.SetDefault(configMaxLimit, defaultMaxLimit)
	viper.SetDefault(configLogLevel, defaultLogLevel)
	viper.SetDefault(configLogFormat, defaultLogFormat)
	viper.SetDefault(configLogFile, defaultLogFile)

	// Bind environment variables with their keys
	checkErr(viper.BindEnv(configFolder, envVarDir))
	checkErr(viper.BindEnv(configAutoCommit, envVarAutoCommit))
	checkErr(viper.BindEnv(configPageSize, envVarPageSize))
	checkErr(viper.BindEnv(configMaxLimit, envVarMaxLimit))
	checkErr(viper.BindEnv(configLogLevel, envVarLogLevel))
	checkErr(viper.BindEnv(configLogFormat, envVarLogFormat))
	checkErr(viper.BindEnv(configLogFile, envVarLogFile))
}

// GetDefaultPageSize returns the configured default page size
func GetDefaultPageSize() int {
	return viper.GetInt(configPageSize)
}

// GetMaxLimit returns the configured maximum limit for pagination
func GetMaxLimit() int {
	return viper.GetInt(configMaxLimit)
}

// ApplyDefaultPagination applies default pagination values if not set
func ApplyDefaultPagination(limit, offset *int) (*int, *int) {
	// Don't apply defaults if user explicitly set offset (means they want pagination)
	if offset != nil && *offset > 0 {
		return limit, offset
	}

	// Don't apply defaults if user explicitly set limit
	if limit != nil && *limit > 0 {
		// Enforce max limit
		maxLimit := GetMaxLimit()
		if *limit > maxLimit {
			enforcedLimit := maxLimit
			return &enforcedLimit, offset
		}
		return limit, offset
	}

	// No pagination requested by user
	return limit, offset
}

func checkErr(err error) {
	if err != nil {
		logging.Error("binding environment variables", "err", err)
	}
}

func setRootPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(configFolder, defaultFolder, "Directory for backlog tasks")
	cmd.PersistentFlags().Bool(configAutoCommit, defaultAutoCommit, "Auto-committing changes to git repository")
	cmd.PersistentFlags().Int(configPageSize, defaultPageSize, "Default page size for pagination")
	cmd.PersistentFlags().Int(configMaxLimit, defaultMaxLimit, "Maximum limit for pagination")
	cmd.PersistentFlags().String(configLogLevel, defaultLogLevel, "Log level (debug, info, warn, error)")
	cmd.PersistentFlags().String(configLogFormat, defaultLogFormat, "Log format (json, text)")
	cmd.PersistentFlags().String(configLogFile, defaultLogFile, "Log file path (defaults to stderr)")

	// Bind flags to viper
	checkErr(viper.BindPFlag(configFolder, cmd.PersistentFlags().Lookup(configFolder)))
	checkErr(viper.BindPFlag(configAutoCommit, cmd.PersistentFlags().Lookup(configAutoCommit)))
	checkErr(viper.BindPFlag(configPageSize, cmd.PersistentFlags().Lookup(configPageSize)))
	checkErr(viper.BindPFlag(configMaxLimit, cmd.PersistentFlags().Lookup(configMaxLimit)))
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
