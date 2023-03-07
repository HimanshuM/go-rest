func {{ .Handler }}(g *gin.Context) {
    var err error
{{- if .Param }}
    {{ .Param }} := g.Param("{{ .Param }}")
{{ end -}}
{{- if .Request }}
    {{ .Request.Name }} := {{ .Request.TypeDeclaration }}{}
    if err = g.ShouldBindJSON({{ .Request.Name }}); err != nil {
        for _, fieldErr := range err.(validator.ValidationErrors) {
            g.AbortWithStatusJSON(http.StatusUnprocessableEntity, fmt.Sprint(fieldErr))
            return
        }
    }
{{- end -}}
{{- if .Response }}
    var {{ .Response.Name }} {{ .Response.TypeDeclaration }}
{{- end }}
    if {{ .Returns }} = {{ .LevelServerHandler }}.{{ .Method }}({{ .Params }}); err != nil {
        g.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
        return
    }
{{- if .Response }}
    g.JSON({{ .HTTPCode }}, &{{ .Response.Name }})
{{ else }}
    g.JSON({{ .HTTPCode }}, gin.H{"status": "Ok"})
{{ end -}}
}
