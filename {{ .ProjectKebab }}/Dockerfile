{{- if .Computed.feature_docker -}}
FROM gcr.io/distroless/static
ARG TARGETPLATFORM
ENTRYPOINT ["/{{ .Scaffold.gomod | pathBase }}"]
COPY ${TARGETPLATFORM}/{{ .Scaffold.gomod | pathBase }} /
{{- end }}
