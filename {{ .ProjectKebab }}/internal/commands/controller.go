// Package commands contains the CLI commands for the application
package commands

import (
	"fmt"

	"context"
)

type Flags struct {
	LogLevel         string
}

type Controller struct {
	Flags *Flags
}

{{ range .Scaffold.commands }}
func (c *Controller) {{ . | toTitleCase | replace "-" "" }}(ctx context.Context) error {
	fmt.Println("Hello World!")
	return nil
}
{{ end }}
