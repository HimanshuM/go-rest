func {{ .Handler }}(g *gin.Context) {
{{- if .URLParam }}
    {{ .URLParam.Name }} = g.Param("{{ .URLParam.Name }}")
{{- end }}
{{- if .JSONRequest }}
    {{ .JSONRequest.Name }} := &{{ .JSONRequest.Type }}{}
    err := g.ShouldBindJSON(&{{ .JSONRequest.Name }})
    if err != nil {
    }
{{- end }}
}
