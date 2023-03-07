func (s *{{ .LevelServer }}) {{ .Handler }}({{ .Params }}) {{ .ReturnType }} {
{{- if .Param }}
    // obj, err := struct.Find({{ .Param }})
    // if err != nil {
    {{- if .Response }}
    //    return {{ .Response.TypeDeclaration }}{}, nil
    {{ else }}
    //    return nil
    {{ end -}}
    // }
{{ end -}}
{{- if .Response }}
    return {{ .Response.TypeDeclaration }}{}, nil
{{ else }}
    return nil
{{ end -}}
}
