{{- if .Computed.feature_docker_integration -}}
//go:build integration

package integration

import (
	"strings"
	"testing"
)

func TestSmokeHelp(t *testing.T) {
	h := NewHarness(t)
	out, err := h.Run("help")
	if err != nil {
		t.Fatalf("help command failed: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "{{ .Project }}") {
		t.Fatalf("expected help output to mention {{ .Project }}, got:\n%s", out)
	}
}

func TestSmokeVersion(t *testing.T) {
	h := NewHarness(t)
	out, err := h.Run("--version")
	if err != nil {
		t.Fatalf("version command failed: %v\noutput: %s", err, out)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected version output, got empty string")
	}
}

func TestSmokeCommand(t *testing.T) {
	h := NewHarness(t)
	out, err := h.Run("{{ index .Scaffold.commands 0 }}")
	if err != nil {
		t.Fatalf("command failed: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Hello World!") {
		t.Fatalf("expected command output to contain Hello World!, got:\n%s", out)
	}
}
{{- end }}
