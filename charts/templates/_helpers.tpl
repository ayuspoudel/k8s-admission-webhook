{{- define "webhook-server.name" -}}
webhook-server
{{- end }}

{{- define "webhook-server.fullname" -}}
{{ include "webhook-server.name" . }}
{{- end }}
