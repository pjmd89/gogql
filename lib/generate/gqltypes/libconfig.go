package gqltypes

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
)

func LibConfigtmpl(types generate.RenderTypes) {
	mt, err := template.New("libconfig.tmpl").Parse(string(generate.LibConfigtmpl))
	if err != nil {
		panic(err)
	}
	if len(types.ScalarType) > 0 {
		types.IsScalar = true
	}
	dir := filepath.Dir(types.LibConfigPath)
	os.MkdirAll(dir, 0770)
	modelFile, err := os.Create(types.LibConfigPath)
	if err != nil {
		panic(err)
	}

	var tmpl bytes.Buffer
	err = mt.Execute(&tmpl, nil)
	if err != nil {
		panic(err)
	}
	x, err := format.Source(tmpl.Bytes())
	if err != nil {
		log.Fatal(err.Error())
	}
	modelFile.Write(x)
}
