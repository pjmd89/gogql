package gqltypes

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
)

func Maintmpl(types generate.RenderTypes) {
	mt, err := template.New("main.tmpl").Parse(string(generate.Maintmpl))
	if err != nil {
		panic(err)
	}
	if len(types.ScalarType) > 0 {
		types.IsScalar = true
	}
	dir := filepath.Dir(types.MainPath)
	os.MkdirAll(dir, 0770)
	modelFile, err := os.Create(types.MainPath)
	if err != nil {
		panic(err)
	}

	var tmpl bytes.Buffer
	err = mt.Execute(&tmpl, types)
	if err != nil {
		panic(err)
	}
	/*
		x, err := format.Source(tmpl.Bytes())
		if err != nil {
			log.Fatal(err.Error())
		}
		modelFile.Write(x)
	*/
	modelFile.Write(tmpl.Bytes())
}
