import (
{{- range $pkg := . }}
    "{{ $pkg }}"
{{- end }}
)
