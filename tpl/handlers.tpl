package {{ .Package }}

{{ .Imports }}
type {{ .LevelServer }} struct{}
{{ range $fn := .Methods }}
{{ $fn }}
{{- end }}
func Use{{ .LevelServer }}() {
    {{ .RoutesPackage }}.Use{{ .LevelServer }}Handler(&{{ .LevelServer }}{})
}
