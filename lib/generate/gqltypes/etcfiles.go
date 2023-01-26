package gqltypes

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
)

func EtcFilestmpl(types generate.RenderTypes) {
	configfiles(types.ConfigJSON, generate.ConfigJSONtmpl)
	configfiles(types.ConfigDB, generate.ConfigDBtmpl)
	configfiles(types.ConfigHTTP, generate.ConfigHTTPtmpl)
}
func configfiles(path string, tmplstr []byte) {
	mt, err := template.New(path).Parse(string(tmplstr))
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0770)
	modelFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	var tmpl bytes.Buffer
	err = mt.Execute(&tmpl, nil)
	if err != nil {
		panic(err)
	}
	modelFile.Write(tmpl.Bytes())
}
