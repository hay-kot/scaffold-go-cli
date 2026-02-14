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

// {{ .Each.Item | toTitleCase | replace "-" "" }}Cmd implements the {{ toKebabCase .Each.Item }} command
type {{ .Each.Item | toTitleCase | replace "-" "" }}Cmd struct {
	flags *Flags
}

// New{{ .Each.Item | toTitleCase | replace "-" "" }}Cmd creates a new {{ toKebabCase .Each.Item }} command
func New{{ .Each.Item | toTitleCase | replace "-" "" }}Cmd(flags *Flags) *{{ .Each.Item | toTitleCase | replace "-" "" }}Cmd {
	return &{{ .Each.Item | toTitleCase | replace "-" "" }}Cmd{flags: flags}
}

// Register adds the {{ toKebabCase .Each.Item }} command to the application
func (cmd *{{ .Each.Item | toTitleCase | replace "-" "" }}Cmd) Register(app *cli.Command) *cli.Command {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "{{ toKebabCase .Each.Item }}",
		Usage: "{{ toKebabCase .Each.Item }} command",
		Flags: []cli.Flag{
			// Add command-specific flags here
		},
		Action: cmd.run,
	})

	return app
}

func (cmd *{{ .Each.Item | toTitleCase | replace "-" "" }}Cmd) run(ctx context.Context, c *cli.Command) error {
	log.Info().Msg("running {{ toKebabCase .Each.Item }} command")

{{- if .Scaffold.feature_json_output }}
	if cmd.flags.JSON {
		return json.NewEncoder(os.Stdout).Encode(map[string]string{
			"message": "Hello World!",
		})
	}
{{- end }}

	fmt.Println("Hello World!")

	return nil
}
