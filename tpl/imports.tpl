import (
{{- range $pkg := . }}
    {{ if ne $pkg.Alias $pkg.Name }}{{ print $pkg.Alias " " }}{{ end }}"{{ $pkg.Path }}"
{{- end }}
)
