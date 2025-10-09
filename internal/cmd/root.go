package cmd

import (
	"context"
	"os"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
	mcpserver "github.com/veggiemonk/backlog/internal/mcp"
	"github.com/veggiemonk/backlog/internal/paths"
)

const (
	envPrefix        = "BACKLOG"
	envVarLogFile    = envPrefix + "_LOG_FILE"
	envVarLogLevel   = envPrefix + "_LOG_LEVEL"
	envVarLogFormat  = envPrefix + "_LOG_FORMAT"
	envVarAutoCommit = envPrefix + "_AUTO_COMMIT"
	envVarDir        = envPrefix + "_FOLDER"

	configFolder      = "folder"
	defaultFolder     = ".backlog"
	configAutoCommit  = "auto-commit"
	defaultAutoCommit = false
	configLogLevel    = "log-level"
	defaultLogLevel   = "info"
	configLogFormat   = "log-format"
	defaultLogFormat  = "text"
	configLogFile     = "log-file"
	defaultLogFile    = ""
)

type runtime struct {
	store      mcpserver.TaskStore
	tasksDir   string
	autoCommit bool
}

type config struct {
	fs          afero.Fs
	store       mcpserver.TaskStore
	tasksDir    string
	skipLogging bool
}

// Option configures the root command construction.
type Option func(*config)

// WithFilesystem injects the filesystem to use when the store is created.
func WithFilesystem(fs afero.Fs) Option {
	return func(cfg *config) {
		cfg.fs = fs
	}
}

// WithStore injects an already-configured task store for the CLI.
func WithStore(store mcpserver.TaskStore) Option {
	return func(cfg *config) {
		cfg.store = store
	}
}

// WithTasksDir preconfigures the tasks directory used when creating a store.
func WithTasksDir(dir string) Option {
	return func(cfg *config) {
		cfg.tasksDir = dir
	}
}

// WithSkipLogging disables logging initialisation and teardown.
func WithSkipLogging(skip bool) Option {
	return func(cfg *config) {
		cfg.skipLogging = skip
	}
}

// NewCommand constructs the backlog root command and its sub-commands.
func NewCommand(opts ...Option) *cli.Command {
	cfg := config{fs: afero.NewOsFs()}
	for _, opt := range opts {
		opt(&cfg)
	}

	rt := &runtime{}
	if cfg.store != nil {
		rt.store = cfg.store
	}
	if cfg.tasksDir != "" {
		rt.tasksDir = cfg.tasksDir
	}

	root := &cli.Command{
		Name:  "backlog",
		Usage: "Backlog is a git-native, markdown-based task manager",
		Description: `A Git-native, Markdown-based task manager for developers and AI agents.
Backlog helps you manage tasks within your git repository.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    configFolder,
				Usage:   "Directory for backlog tasks",
				Value:   defaultFolder,
				Sources: cli.EnvVars(envVarDir),
			},
			&cli.BoolFlag{
				Name:    configAutoCommit,
				Usage:   "Auto-committing changes to git repository",
				Value:   defaultAutoCommit,
				Sources: cli.EnvVars(envVarAutoCommit),
			},
			&cli.StringFlag{
				Name:    configLogLevel,
				Usage:   "Log level (debug, info, warn, error)",
				Value:   defaultLogLevel,
				Sources: cli.EnvVars(envVarLogLevel),
			},
			&cli.StringFlag{
				Name:    configLogFormat,
				Usage:   "Log format (json, text)",
				Value:   defaultLogFormat,
				Sources: cli.EnvVars(envVarLogFormat),
			},
			&cli.StringFlag{
				Name:    configLogFile,
				Usage:   "Log file path (defaults to stderr)",
				Value:   defaultLogFile,
				Sources: cli.EnvVars(envVarLogFile),
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if !cfg.skipLogging {
				logging.Init(
					cmd.String(configLogLevel),
					cmd.String(configLogFormat),
					cmd.String(configLogFile),
				)
			}

			rt.autoCommit = cmd.Bool(configAutoCommit)

			if rt.store != nil {
				if rt.tasksDir == "" {
					rt.tasksDir = cmd.String(configFolder)
				}
				return ctx, nil
			}

			tasksDir := cmd.String(configFolder)
			logging.Debug("resolve env var", configFolder, tasksDir, configAutoCommit, rt.autoCommit)

			resolvedDir, err := paths.ResolveTasksDir(cfg.fs, tasksDir)
			if err != nil {
				logging.Error("tasks directory", "error", err)
				resolvedDir = tasksDir
			}
			logging.Debug("resolve tasks directory", configFolder, resolvedDir)

			rt.tasksDir = resolvedDir
			rt.store = core.NewFileTaskStore(cfg.fs, resolvedDir)

			return ctx, nil
		},
		After: func(ctx context.Context, _ *cli.Command) error {
			if !cfg.skipLogging {
				logging.Close()
			}
			return nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := cli.ShowCommandHelp(ctx, cmd, ""); err != nil {
				return err
			}
			return nil
		},
	}

	root.Commands = []*cli.Command{
		newArchiveCommand(rt),
		newCreateCommand(rt),
		newDoctorCommand(rt),
		newEditCommand(rt),
		newInstructionsCommand(),
		newListCommand(rt),
		newMCPCommand(rt),
		newVersionCommand(),
		newViewCommand(rt),
	}

	return root
}

// Execute runs the backlog CLI with the given arguments.
func Execute() {
	cmd := NewCommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logging.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
