package main

import (
	"fmt"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"{{ .Scaffold.gomod }}/app/commands"
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
	ctrl := &commands.Controller{}

	app := &cli.App{
		Name:    "{{ .Project }}",
		Usage:   "{{ .Scaffold.description }}",
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cwd",
				Usage: "current working directory",
				Value: ".",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "log level (debug, info, warn, error, fatal, panic)",
				Value: "panic",
			},
		},
		Before: func(ctx *cli.Context) error {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

			ctrl.WorkingDir = ctx.String("cwd")

			switch ctx.String("log-level") {
			case "debug":
				log.Level(zerolog.DebugLevel)
			case "info":
				log.Level(zerolog.InfoLevel)
			case "warn":
				log.Level(zerolog.WarnLevel)
			case "error":
				log.Level(zerolog.ErrorLevel)
			case "fatal":
				log.Level(zerolog.FatalLevel)
			default:
				log.Level(zerolog.PanicLevel)
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "hello",
				Usage:  "Says hello world",
				Action: ctrl.HelloWorld,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run {{ .Project }}")
	}
}
