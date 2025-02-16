package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"{{ .Scaffold.gomod }}/internal/commands"
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

func main() {
	ctrl := &commands.Controller{
		Flags: &commands.Flags{},
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	app := &cli.Command{
		Name:    "{{ .Project }}",
		Usage:   `{{ .Scaffold.description }}`,
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level (debug, info, warn, error, fatal, panic)",
				Sources: cli.EnvVars("LOG_FORMAT"),
				Value:       "panic",
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			level, err := zerolog.ParseLevel(c.String("log-level"))
			if err != nil {
				return ctx, fmt.Errorf("failed to parse log level: %w", err)
			}

			log.Logger = log.Level(level)

			return ctx, nil
		},
		Commands: []*cli.Command{
      {{-  range .Scaffold.commands }}
			{
				Name:   "{{ toKebabCase . }}",
				Usage:  "",
				Action: func(ctx context.Context, c *cli.Command) error {
					return ctrl.{{ . | toTitleCase | replace "-" "" }}(ctx)
				},
			},
			{{- end }}
		},
	}

	ctx := context.Background()

	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run {{ .Project }}")
	}
}
