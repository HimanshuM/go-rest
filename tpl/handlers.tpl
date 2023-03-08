package {{ .Package }}

{{ .Imports -}}
{{ if .HasLevelServer }}
type {{ .LevelServer }} struct{}
{{ end -}}
{{ range $fn := .Methods }}
{{ $fn }}
{{- end }}
{{- if .HasLevelServer }}
func use{{ .LevelServer }}() {
    {{ .RoutesPackage }}.Use{{ .LevelServer }}Handler(&{{ .LevelServer }}{})
}
{{ end }}
func {{ .InitName }}{{ .Level }}() {
{{- if .HasLevelServer }}
    use{{ .LevelServer }}()
{{- end }}
{{- range $subPath := .SubPaths }}
    {{ $subPath }}()
{{- end }}
}
