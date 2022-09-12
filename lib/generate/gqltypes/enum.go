package gqltypes

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gqlparser/v2/ast"
)

func NewEnum(render generate.GqlGenerate, key string, value *ast.Definition) generate.EnumDef {
	structDef := generate.EnumDef{}
	structDef.Attr = make([]generate.EnumAttrDef, 0)
	structDef.Name = strings.Title(key)
	structDef.PackageName = render.EnumPath
	for _, vValue := range value.EnumValues {
		enumAttrDef := generate.EnumAttrDef{}
		enumAttrDef.Name = strings.ToUpper(key + "_" + vValue.Name)
		enumAttrDef.Value = vValue.Name
		structDef.Attr = append(structDef.Attr, enumAttrDef)
	}
	structDef.FilePath = render.ModulePath + "/generate/" + render.ModelPath + "/" + render.EnumPath + "/" + strings.ToLower(key) + ".go"

	return structDef
}
func EnumTmpl(types generate.RenderTypes) {
	ut, err := template.New("enum.tmpl").Parse(string(generate.Enumtmpl))
	if err != nil {
		panic(err)
	}
	for _, v := range types.EnumType {
		dir := filepath.Dir(v.FilePath)
		os.MkdirAll(dir, 0770)
		modelFile, err := os.Create(v.FilePath)
		if err != nil {
			panic(err)
		}

		var tmpl bytes.Buffer
		err = ut.Execute(&tmpl, v)
		if err != nil {
			panic(err)
		}
		x, _ := format.Source(tmpl.Bytes())
		modelFile.Write(x)
	}
}
