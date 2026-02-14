{{- if .Computed.feature_docker -}}
FROM gcr.io/distroless/static
ENTRYPOINT ["/{{ .Scaffold.gomod | pathBase }}"]
COPY {{ .Scaffold.gomod | pathBase }} /
{{- end }}
