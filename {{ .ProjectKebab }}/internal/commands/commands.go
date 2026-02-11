package commands

import (
	"context"
{{- if .Scaffold.feature_json_output }}
	"encoding/json"
{{- end }}
	"fmt"
{{- if .Scaffold.feature_json_output }}
	"os"
{{- end }}

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

{{- range .Scaffold.commands }}

// {{ . | toTitleCase | replace "-" "" }}Cmd implements the {{ toKebabCase . }} command
type {{ . | toTitleCase | replace "-" "" }}Cmd struct {
	flags *Flags
}

// New{{ . | toTitleCase | replace "-" "" }}Cmd creates a new {{ toKebabCase . }} command
func New{{ . | toTitleCase | replace "-" "" }}Cmd(flags *Flags) *{{ . | toTitleCase | replace "-" "" }}Cmd {
	return &{{ . | toTitleCase | replace "-" "" }}Cmd{flags: flags}
}

// Register adds the {{ toKebabCase . }} command to the application
func (cmd *{{ . | toTitleCase | replace "-" "" }}Cmd) Register(app *cli.Command) *cli.Command {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "{{ toKebabCase . }}",
		Usage: "{{ toKebabCase . }} command",
		Flags: []cli.Flag{
			// Add command-specific flags here
		},
		Action: cmd.run,
	})

	return app
}

func (cmd *{{ . | toTitleCase | replace "-" "" }}Cmd) run(ctx context.Context, c *cli.Command) error {
	log.Info().Msg("running {{ toKebabCase . }} command")

{{- if $.Scaffold.feature_json_output }}
	if cmd.flags.JSON {
		return json.NewEncoder(os.Stdout).Encode(map[string]string{
			"message": "Hello World!",
		})
	}
{{- end }}

	fmt.Println("Hello World!")

	return nil
}
{{- end }}
