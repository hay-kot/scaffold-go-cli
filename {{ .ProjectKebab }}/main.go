package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
{{- if .Scaffold.feature_file_logging }}
	"io"
	"path/filepath"
{{- end }}

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"{{ .Scaffold.gomod }}/internal/commands"
{{- if .Scaffold.feature_config_file }}
	"{{ .Scaffold.gomod }}/internal/config"
{{- end }}
{{- if .Scaffold.feature_file_logging }}
	"{{ .Scaffold.gomod }}/internal/paths"
{{- end }}
)

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

{{- if .Scaffold.feature_file_logging }}
func setupLogger(level string, logFile string) error {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	var output io.Writer = zerolog.ConsoleWriter{Out: os.Stderr}

	if logFile != "" {
		// Create log directory if it doesn't exist
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open log file
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}

		// Write to both console and file
		output = io.MultiWriter(
			zerolog.ConsoleWriter{Out: os.Stderr},
			file,
		)
	}

	log.Logger = log.Output(output).Level(parsedLevel)

	return nil
}
{{- else }}
func setupLogger(level string) error {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(parsedLevel)

	return nil
}
{{- end }}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flags := &commands.Flags{}

	app := &cli.Command{
		Name:                  "{{ .Project }}",
		Usage:                 `{{ .Scaffold.description }}`,
		Version:               build(),
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level (debug, info, warn, error, fatal, panic)",
				Sources:     cli.EnvVars("LOG_LEVEL"),
				Value:       "info",
				Destination: &flags.LogLevel,
			},
{{- if .Scaffold.feature_file_logging }}
			&cli.StringFlag{
				Name:        "log-file",
				Usage:       "path to log file (optional)",
				Sources:     cli.EnvVars("LOG_FILE"),
				Destination: &flags.LogFile,
			},
{{- end }}
{{- if .Scaffold.feature_config_file }}
			&cli.StringFlag{
				Name:        "config",
				Usage:       "path to config file",
				Sources:     cli.EnvVars("CONFIG_FILE"),
				Destination: &flags.ConfigFile,
			},
{{- end }}
{{- if .Scaffold.feature_json_output }}
			&cli.BoolFlag{
				Name:        "json",
				Usage:       "output in JSON format",
				Destination: &flags.JSON,
			},
{{- end }}
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
{{- if .Scaffold.feature_config_file }}
			cfg, err := func() (config.Config, error) {
				if flags.ConfigFile != "" {
					return config.ReadFrom(flags.ConfigFile)
				}
				return config.Read()
			}()
			if err != nil {
				return ctx, fmt.Errorf("loading config: %w", err)
			}

			if flags.LogLevel == "info" && cfg.LogLevel != "" {
				flags.LogLevel = cfg.LogLevel
			}
{{- if .Scaffold.feature_file_logging }}
			if flags.LogFile == "" && cfg.LogFile != "" {
				flags.LogFile = cfg.LogFile
			}
{{- end }}
{{- end }}
{{- if .Scaffold.feature_file_logging }}
			logFile := flags.LogFile
			if logFile == "" {
				logFile = filepath.Join(paths.DataDir(), "{{ .Scaffold.gomod | pathBase }}.log")
			}

			if err := setupLogger(flags.LogLevel, logFile); err != nil {
				return ctx, err
			}
{{- else }}
			if err := setupLogger(flags.LogLevel); err != nil {
				return ctx, err
			}
{{- end }}

			return ctx, nil
		},
	}

{{- range .Scaffold.commands }}
	app = commands.New{{ . | toTitleCase | replace "-" "" }}Cmd(flags).Register(app)
{{- end }}
	// +scaffold:command:register

	exitCode := 0
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, os.Args); err != nil {
		const colorRed = "\033[38;2;215;95;107m"
		const colorGray = "\033[38;2;163;163;163m"
		const colorReset = "\033[0m"
		fmt.Fprintf(os.Stderr, "\n%s╭ Error%s\n%s│%s %s%s%s\n%s╵%s\n",
			colorRed, colorReset,
			colorRed, colorReset, colorGray, err.Error(), colorReset,
			colorRed, colorReset,
		)
		exitCode = 1
	}

	os.Exit(exitCode)
}
