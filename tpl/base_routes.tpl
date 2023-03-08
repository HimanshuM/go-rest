package {{ .Package }}

{{ .Imports }}
func {{ .SetupCase }}etup{{ .Level }}Routes(server *gin.RouterGroup) {
{{- if .Server }}
    {{ .Server }} := server.Group("{{ .Route }}")
{{- end -}}
{{- range $line := .Lines }}
    {{ $line }}
{{- end }}
}
