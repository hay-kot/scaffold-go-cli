package commands

// Flags holds global flags shared across all commands
type Flags struct {
	LogLevel string
{{- if .Scaffold.feature_file_logging }}
	LogFile  string
{{- end }}
}
