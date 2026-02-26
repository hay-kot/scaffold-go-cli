{{- if .Computed.feature_docker_integration -}}
//go:build integration

package integration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	imageName = "{{ .Scaffold.gomod | pathBase }}-integration"
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Fprintf(os.Stderr, "setup failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	if err := teardown(); err != nil {
		fmt.Fprintf(os.Stderr, "teardown failed: %v\n", err)
	}

	os.Exit(code)
}

func setup() error {
	// Build the binary with goreleaser snapshot
	cmd := exec.Command("goreleaser", "build", "--snapshot", "--clean", "--single-target")
	cmd.Dir = ".."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("goreleaser build: %w", err)
	}

	// Find the built binary
	binary, err := findBinary()
	if err != nil {
		return fmt.Errorf("find binary: %w", err)
	}

	// Copy binary into integration directory for docker build context
	input, err := os.ReadFile(binary)
	if err != nil {
		return fmt.Errorf("read binary: %w", err)
	}

	if err := os.WriteFile("{{ .Scaffold.gomod | pathBase }}", input, 0o755); err != nil {
		return fmt.Errorf("write binary: %w", err)
	}

	// Build docker image
	cmd = exec.Command("docker", "build", "-t", imageName, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build: %w", err)
	}

	return nil
}

func teardown() error {
	// Remove copied binary
	os.Remove("{{ .Scaffold.gomod | pathBase }}")

	// Remove docker image
	cmd := exec.Command("docker", "rmi", "-f", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findBinary() (string, error) {
	// goreleaser places binaries in dist/<project>_<os>_<arch>/
	entries, err := os.ReadDir("../dist")
	if err != nil {
		return "", fmt.Errorf("read dist: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := fmt.Sprintf("../dist/%s/{{ .Scaffold.gomod | pathBase }}", entry.Name())
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no binary found in dist/")
}

// runCLI executes the CLI binary inside the docker container with the given
// arguments and returns the combined stdout output.
func runCLI(args ...string) (string, error) {
	cmdArgs := append([]string{"run", "--rm", imageName}, args...)
	cmd := exec.Command("docker", cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker run failed: %w\nstderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func TestVersion(t *testing.T) {
	out, err := runCLI("--version")
	if err != nil {
		t.Fatalf("failed to get version: %v", err)
	}

	if out == "" {
		t.Fatal("expected version output, got empty string")
	}

	t.Logf("version output: %s", out)
}
{{- end }}
