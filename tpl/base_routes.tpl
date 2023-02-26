package {{ .Package }}

{{ .Imports }}
func Setup{{ .Level }}Routes(server *gin.Engine) {
{{- range $line := .Lines }}
    {{ $line }}
{{- end }}
}
