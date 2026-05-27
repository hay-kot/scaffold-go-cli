# {{ .Project }}

{{ .Scaffold.description }}

## Installation

```bash
go install {{ .Scaffold.gomod }}
```

## Usage

TODO
{{- if .Computed.feature_profiling }}

## Profiling

This project can expose Go pprof endpoints or write profiles to disk when profiling is enabled.

Expose the pprof HTTP endpoint:

```bash
{{ .Project }} --pprof <command>
```

By default, this listens on `127.0.0.1:{{ .Computed.pprof_port }}`. Override it with `--pprof-addr` if needed:

```bash
{{ .Project }} --pprof-addr 127.0.0.1:6060 <command>
```

Then open <http://127.0.0.1:{{ .Computed.pprof_port }}/debug/pprof/> or collect a CPU profile:

```bash
go tool pprof http://127.0.0.1:{{ .Computed.pprof_port }}/debug/pprof/profile?seconds=30
```

Write profiles to disk for short-lived CLI commands:

```bash
{{ .Project }} --cpu-profile cpu.pprof --heap-profile heap.pprof <command>
go tool pprof cpu.pprof
go tool pprof heap.pprof
```
{{- end }}
