package gqltypes

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gqlparser/v2/ast"
)

func NewUnion(key string, value *ast.Definition) generate.UnionAttrDef {
	unionDef := generate.UnionAttrDef{}
	unionDef.Name = strings.Title(key)
	types := make([]string, 0)
	for _, vValue := range value.Types {
		types = append(types, strings.Title(vValue))
	}
	unionDef.Types = strings.Join(types, " | ")
	return unionDef
}
func UnionTmpl(types generate.RenderTypes) {
	et, err := template.New("union.tmpl").Parse(string(generate.Uniontmpl))
	if err != nil {
		panic(err)
	}
	if types.UnionType.FilePath != "" {
		dir := filepath.Dir(types.UnionType.FilePath)
		os.MkdirAll(dir, 0770)
		modelFile, err := os.Create(types.UnionType.FilePath)
		if err != nil {
			panic(err.Error() + " - " + types.UnionType.FilePath)
		}

		var tmpl bytes.Buffer
		err = et.Execute(&tmpl, types.UnionType)
		if err != nil {
			fmt.Printf("error en union %s: %s", types.UnionType.FilePath, err.Error())
			panic(err)
		}
		x, _ := format.Source(tmpl.Bytes())
		modelFile.Write(x)
	}
}
