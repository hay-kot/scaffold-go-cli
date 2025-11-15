package main

import (
	"context"
	"fmt"
	"os"
{{- if .Scaffold.feature_file_logging }}
	"io"
	"path/filepath"
{{- end }}

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"{{ .Scaffold.gomod }}/internal/commands"
	"{{ .Scaffold.gomod }}/internal/printer"
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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	log.Logger = log.Level(parsedLevel)

	return nil
}
{{- end }}

func main() {
{{- if .Scaffold.feature_file_logging }}
	if err := setupLogger("info", ""); err != nil {
		panic(err)
	}
{{- else }}
	if err := setupLogger("info"); err != nil {
		panic(err)
	}
{{- end }}

	p := printer.New(os.Stderr)
	ctx := printer.NewContext(context.Background(), p)

	flags := &commands.Flags{}

	app := &cli.Command{
		Name:    "{{ .Project }}",
		Usage:   `{{ .Scaffold.description }}`,
		Version: build(),
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
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
{{- if .Scaffold.feature_file_logging }}
			if err := setupLogger(flags.LogLevel, flags.LogFile); err != nil {
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

	exitCode := 0
	if err := app.Run(ctx, os.Args); err != nil {
		fmt.Println()
		printer.Ctx(ctx).FatalError(err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
