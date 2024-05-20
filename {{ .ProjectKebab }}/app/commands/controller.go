// Package commands contains the CLI commands for the application
package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type Flags struct {
	LogLevel         string
}

type Controller struct {
	Flags *Flags
}

{{ range .Scaffold.commands }}
func (c *Controller) {{ . | titlecase | replace "-" "" }}(ctx *cli.Context) error {
	fmt.Println("Hello World!")
	return nil
}
{{ end }}
