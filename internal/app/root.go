// Package app sets up the CLI application using Cobra.
package app

import (
	"fmt"

	"github.com/spf13/cobra"
	appcrypto "github.com/xef/xefcli/internal/app/crypto"
	"github.com/xef/xefcli/internal/app/dev"
	"github.com/xef/xefcli/internal/app/file"
	"github.com/xef/xefcli/internal/app/git"
	apphttp "github.com/xef/xefcli/internal/app/http"
	"github.com/xef/xefcli/internal/app/json"
	"github.com/xef/xefcli/internal/app/system"
	"github.com/xef/xefcli/internal/core/config"
	"github.com/xef/xefcli/internal/core/logger"
	infracrypto "github.com/xef/xefcli/internal/infrastructure/crypto"
	"github.com/xef/xefcli/internal/infrastructure/filesystem"
	"github.com/xef/xefcli/internal/infrastructure/network"
	"github.com/xef/xefcli/internal/infrastructure/systeminfo"
	"github.com/xef/xefcli/internal/pkg/tui"
)

// App holds application dependencies.
type App struct {
	RootCmd *cobra.Command
	Config  *config.Loader
	Logger  logger.Logger
}

// New creates a new App instance with all dependencies wired.
func New(version string) *App {
	cfg := config.NewLoader()
	log, _ := logger.New(logger.Config{Level: logger.InfoLevel, Format: "pretty"})

	app := &App{
		Config: cfg,
		Logger: log,
	}

	app.RootCmd = &cobra.Command{
		Use:   "xef",
		Short: "XefCLI - The Ultimate Developer Toolkit",
		Long: `XefCLI is a production-grade, cross-platform CLI toolkit for developers.

It provides powerful tools for file management, data processing, cryptography,
HTTP operations, system monitoring, and development workflows.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfgFile, _ := cmd.Flags().GetString("config")
			if err := cfg.Load(cfgFile); err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			level, _ := cmd.Flags().GetString("log-level")
			format, _ := cmd.Flags().GetString("log-format")
			newLog, err := logger.New(logger.Config{
				Level:  logger.LogLevel(level),
				Format: format,
			})
			if err == nil {
				app.Logger = newLog
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	app.RootCmd.PersistentFlags().StringP("config", "c", "", "config file path")
	app.RootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")
	app.RootCmd.PersistentFlags().String("log-format", "pretty", "log format (json, pretty)")

	// Initialize infrastructure layer
	fs := filesystem.NewOSFileSystem()
	hasher := infracrypto.NewDefaultHasher()
	httpClient := network.NewDefaultHTTPClient(0)
	sysProvider := systeminfo.NewGopsutilProvider()

	// Register domain commands
	app.RootCmd.AddCommand(file.NewCommand(fs, log))
	app.RootCmd.AddCommand(json.NewCommand(log))
	app.RootCmd.AddCommand(appcrypto.NewCommand(hasher, log))
	app.RootCmd.AddCommand(apphttp.NewCommand(httpClient, log))
	app.RootCmd.AddCommand(git.NewCommand(log))
	app.RootCmd.AddCommand(system.NewCommand(sysProvider, log))
	app.RootCmd.AddCommand(dev.NewCommand(fs, log))

	app.RootCmd.SetUsageTemplate(usageTemplate())
	return app
}

// Execute runs the application.
func (a *App) Execute() error {
	if err := a.RootCmd.Execute(); err != nil {
		tui.PrintError(err.Error())
		return err
	}
	return nil
}

func usageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
