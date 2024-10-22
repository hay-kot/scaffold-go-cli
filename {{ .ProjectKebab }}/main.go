package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

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

	app := &cli.App{
		Name:    "{{ .Project }}",
		Usage:   `{{ .Scaffold.description }}`,
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level (debug, info, warn, error, fatal, panic)",
				Value:       "panic",
			},
		},
		Before: func(ctx *cli.Context) error {
			level, err := zerolog.ParseLevel(ctx.String("log-level"))
			if err != nil {
				return fmt.Errorf("failed to parse log level: %w", err)
			}

			log.Logger = log.Level(level)

			return nil
		},
		Commands: []*cli.Command{
      {{  range .Scaffold.commands -}}
			{
				Name:   "{{ toKebabCase . }}",
				Usage:  "",
				Action: func(ctx *cli.Context) error {
					return ctrl.{{ . | toTitleCase | replace "-" "" }}(ctx.Context)
				},
			},{{end }}
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run {{ .Project }}")
	}
}
