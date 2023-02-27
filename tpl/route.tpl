func {{ .Handler }}(g *gin.Context) {
    var err error
{{- if .Param }}
    {{ .Param }} := g.Param("{{ .Param }}")
{{ end -}}
{{- if .Request }}
    {{ .Request.Name }} := &{{ .Request.Type }}{}
    err = g.ShouldBindJSON(&{{ .Request.Name }})
    if err != nil {
    }
{{- end -}}
{{- if .Response }}
    var {{ .Response.Name }} *{{ .Response.Type }}
{{ end }}
    if {{ .Returns }} = {{ .LevelServerHandler }}.{{ .Method }}({{ .Params }}); err != nil {
    }
{{- if .Response }}
    g.JSON(200, &{{ .Response.Name }})
{{ else }}
    g.JSON(200, g.H{"status": "Ok"})
{{ end -}}
}
