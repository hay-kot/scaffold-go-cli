{{- if .Computed.feature_docker_integration -}}
//go:build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var cliBin string

func TestMain(m *testing.M) {
	if os.Getenv("CLI_INTEGRATION") != "1" {
		fmt.Fprintln(os.Stderr, "skipping integration tests: CLI_INTEGRATION=1 not set (run via mise run integration)")
		os.Exit(0)
	}

	path, err := exec.LookPath("{{ .Scaffold.gomod | pathBase }}")
	if err != nil {
		panic("{{ .Scaffold.gomod | pathBase }} binary not found in PATH; build it first")
	}
	cliBin = path

	os.Exit(m.Run())
}
{{- end }}
