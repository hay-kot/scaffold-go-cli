{{- if .Computed.feature_docker_integration -}}
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY {{ .Scaffold.gomod | pathBase }} /usr/local/bin/{{ .Scaffold.gomod | pathBase }}

ENTRYPOINT ["{{ .Scaffold.gomod | pathBase }}"]
{{- end }}
