package {{ .Package }}

{{ .Imports }}
func {{ .SetupCase }}etup{{ .Level }}Routes(server *gin.RouterGroup) {
{{- if .Server }}
    {{ .Server }} := server.Group("{{ .Route }}")
{{- if .Middlewares }}
{{- range $mw := .Middlewares }}
    {{ $.Server }}.Use({{ $mw }})
{{- end -}}
{{- end -}}
{{- end -}}
{{- range $line := .Lines }}
    {{ $line }}
{{- end }}
}
