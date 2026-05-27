{{- if .Computed.feature_profiling }}
package profiler

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProfilerHTTPServer(t *testing.T) {
	profiler := New(Options{HTTPAddr: "127.0.0.1:0"})

	if err := profiler.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer func() {
		if err := profiler.Stop(shutdownCtx); err != nil {
			t.Fatalf("Stop() error = %v", err)
		}
	}()

	resp, err := http.Get("http://" + profiler.Addr() + "/debug/pprof/")
	if err != nil {
		t.Fatalf("GET pprof index error = %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET pprof index status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestProfilerWritesProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	cpuProfile := filepath.Join(tmpDir, "cpu.pprof")
	heapProfile := filepath.Join(tmpDir, "heap.pprof")

	profiler := New(Options{
		CPUProfile:  cpuProfile,
		HeapProfile: heapProfile,
	})

	if err := profiler.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	for i := 0; i < 100000; i++ {
		_ = i * i
	}

	if err := profiler.Stop(context.Background()); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	assertNonEmptyFile(t, cpuProfile)
	assertNonEmptyFile(t, heapProfile)
}

func assertNonEmptyFile(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if info.Size() == 0 {
		t.Fatalf("%s is empty", path)
	}
}
{{- end }}
