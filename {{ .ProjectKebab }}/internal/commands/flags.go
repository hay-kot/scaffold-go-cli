package commands

// Flags holds global flags shared across all commands
type Flags struct {
	LogLevel string
{{- if .Scaffold.feature_file_logging }}
	LogFile  string
{{- end }}
{{- if .Scaffold.feature_config_file }}
	ConfigFile string
{{- end }}
{{- if .Scaffold.feature_json_output }}
	JSON bool
{{- end }}
}
