func {{ .Handler }}(g *gin.Context) {
    var err error
{{- if .Param }}
    {{ .Param }} := g.Param("{{ .Param }}")
{{ end -}}
{{- if .Request }}
    {{ .Request.Name }} := &{{ .Request.Alias }}.{{ .Request.Type }}{}
    err = g.ShouldBindJSON({{ .Request.Name }})
    if err != nil {
    }
{{- end -}}
{{- if .Response }}
    var {{ .Response.Name }} {{ .Response.TypeDeclaration }}
{{- end }}
    if {{ .Returns }} = {{ .LevelServerHandler }}.{{ .Method }}({{ .Params }}); err != nil {
    }
{{- if .Response }}
    g.JSON({{ .HTTPCode }}, &{{ .Response.Name }})
{{ else }}
    g.JSON({{ .HTTPCode }}, gin.H{"status": "Ok"})
{{ end -}}
}
