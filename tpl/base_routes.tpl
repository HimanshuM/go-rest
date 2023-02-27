package {{ .Package }}

{{ .Imports }}
func Setup{{ .Level }}Routes(server *gin.Engine) {
{{- if .Server }}
    {{ .Server }} := server.Group("{{ .Route }}")
{{- end -}}
{{- range $line := .Lines }}
    {{ $line }}
{{- end }}
}
