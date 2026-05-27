{{- if .Computed.feature_profiling }}
package profiler

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	runtimepprof "runtime/pprof"

	"github.com/rs/zerolog/log"
)

// Options configures runtime profiling.
type Options struct {
	HTTPAddr    string
	CPUProfile  string
	HeapProfile string
}

// Profiler manages optional pprof endpoints and profile dumps.
type Profiler struct {
	options    Options
	server     *http.Server
	listener   net.Listener
	cpuProfile *os.File
}

// New creates a Profiler from options.
func New(options Options) *Profiler {
	return &Profiler{options: options}
}

// Start enables configured profilers.
func (p *Profiler) Start(ctx context.Context) error {
	if p.options.CPUProfile != "" {
		if err := p.startCPUProfile(); err != nil {
			return err
		}
	}

	if p.options.HTTPAddr != "" {
		if err := p.startHTTPServer(); err != nil {
			_ = p.Stop(ctx)
			return err
		}
	}

	return nil
}

// Stop shuts down configured profilers and writes final profile dumps.
func (p *Profiler) Stop(ctx context.Context) error {
	var stopErr error

	if p.cpuProfile != nil {
		runtimepprof.StopCPUProfile()
		if err := p.cpuProfile.Close(); err != nil {
			stopErr = fmt.Errorf("close CPU profile: %w", err)
		}
		p.cpuProfile = nil
		log.Info().Str("path", p.options.CPUProfile).Msg("wrote CPU profile")
	}

	if p.options.HeapProfile != "" {
		if err := writeHeapProfile(p.options.HeapProfile); err != nil {
			if stopErr != nil {
				stopErr = fmt.Errorf("%v; write heap profile: %w", stopErr, err)
			} else {
				stopErr = fmt.Errorf("write heap profile: %w", err)
			}
		} else {
			log.Info().Str("path", p.options.HeapProfile).Msg("wrote heap profile")
		}
	}

	if p.server != nil {
		log.Info().Str("addr", p.Addr()).Msg("shutting down pprof server")
		if err := p.server.Shutdown(ctx); err != nil {
			if stopErr != nil {
				stopErr = fmt.Errorf("%v; shutdown pprof server: %w", stopErr, err)
			} else {
				stopErr = fmt.Errorf("shutdown pprof server: %w", err)
			}
		}
		p.server = nil
	}

	return stopErr
}

// Addr returns the bound HTTP pprof address.
func (p *Profiler) Addr() string {
	if p.listener == nil {
		return ""
	}
	return p.listener.Addr().String()
}

func (p *Profiler) startCPUProfile() error {
	file, err := os.Create(p.options.CPUProfile)
	if err != nil {
		return fmt.Errorf("create CPU profile: %w", err)
	}

	if err := runtimepprof.StartCPUProfile(file); err != nil {
		_ = file.Close()
		return fmt.Errorf("start CPU profile: %w", err)
	}

	p.cpuProfile = file
	log.Info().Str("path", p.options.CPUProfile).Msg("started CPU profiling")
	return nil
}

func (p *Profiler) startHTTPServer() error {
	listener, err := net.Listen("tcp", p.options.HTTPAddr)
	if err != nil {
		return fmt.Errorf("listen for pprof server: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	p.listener = listener
	p.server = &http.Server{Handler: mux}

	go func() {
		if err := p.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("pprof server failed")
		}
	}()

	log.Info().Str("url", fmt.Sprintf("http://%s/debug/pprof/", p.Addr())).Msg("pprof endpoint available")
	return nil
}

func writeHeapProfile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create heap profile: %w", err)
	}
	runtime.GC()
	if err := runtimepprof.WriteHeapProfile(file); err != nil {
		_ = file.Close()
		return fmt.Errorf("write heap profile: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close heap profile: %w", err)
	}

	return nil
}
{{- end }}
