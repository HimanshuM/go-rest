package {{ .Package }}

{{ .Imports }}
type {{ .LevelServer }} interface {
{{- range $def := .Methods }}
    {{ $def }}
{{- end }}
}

var {{ .LevelServerHandler }} {{ .LevelServer }}

func Use{{ .LevelServer }}Handler(handler *{{ .LevelServer }}) {
    {{ .LevelServerHandler }} = *handler
}
{{ range $fn := .Functions }}
{{ $fn }}
{{- end -}}
