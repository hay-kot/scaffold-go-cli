package commands

// Flags holds global flags shared across all commands
type Flags struct {
	LogLevel string
	NoColor  bool
{{- if .Computed.feature_file_logging }}
	LogFile  string
{{- end }}
{{- if .Computed.feature_config_file }}
	ConfigFile string
{{- end }}
{{- if .Computed.feature_json_output }}
	JSON bool
{{- end }}
}
