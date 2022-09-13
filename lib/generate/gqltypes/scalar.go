package gqltypes

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gqlparser/v2/ast"
	"golang.org/x/exp/slices"
)

func NewScalar(render generate.GqlGenerate, key string, value *ast.Definition) (scalarDef *generate.ScalarDef) {
	scalarDef = nil
	if !slices.Contains(generate.OmitScalarTypes, key) {
		typeRegex := regexp.MustCompile(`-type ([^\n]+)`)
		scalarDef = &generate.ScalarDef{}
		scalarDef.Name = strings.Title(key)
		scalarDef.VarName = strings.ToLower(key) + "Scalar"
		scalarDef.PackageName = render.ScalarPath
		typeRegexResult := typeRegex.FindStringSubmatch(value.Description)
		namedTyped := "string"
		if len(typeRegexResult) > 1 {
			namedTyped, _ = generate.GetNamedType(typeRegexResult[1])
		}
		scalarDef.TypeName = namedTyped
		scalarDef.FilePath = render.ModulePath + "/generate/" + render.ResolverPath + "/" + render.ScalarPath + "/" + strings.ToLower(key) + ".go"
	}
	return scalarDef
}
func ScalarTmpl(types generate.RenderTypes) {
	st, err := template.New("scalar.tmpl").Parse(string(generate.Scalartmpl))
	if err != nil {
		panic(err)
	}
	for _, v := range types.ScalarType {
		dir := filepath.Dir(v.FilePath)
		os.MkdirAll(dir, 0770)
		modelFile, err := os.Create(v.FilePath)
		if err != nil {
			panic(err)
		}

		var tmpl bytes.Buffer
		err = st.Execute(&tmpl, v)
		if err != nil {
			panic(err)
		}
		x, _ := format.Source(tmpl.Bytes())
		modelFile.Write(x)
	}
}
