package authorization

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pjmd89/gogql/lib/generate"
	"golang.org/x/exp/slices"
)

var (
	skipObject = []string{"__Type", "__type", "__EnumValue", "__Directive", "__InputValue", "__Schema", "__schema", "__Field"}
)

type AuthData struct {
	TypeNames map[string]string
	Resolvers map[string]string
	Grants    map[string]map[string]map[string][]string
}

func Generate(gqlGenerate generate.GqlGenerate) {

	fmt.Println("generaing auth")

	authRegex := regexp.MustCompile(`-auth ([^\n]+)`)
	authData := AuthData{
		TypeNames: map[string]string{},
		Resolvers: map[string]string{},
		Grants:    map[string]map[string]map[string][]string{},
	}

	for _, v := range gqlGenerate.Schema.Types {
		switch v.Kind {
		case "OBJECT":
			if !slices.Contains(skipObject, strings.Trim(v.Name, " ")) {

				for _, vv := range v.Fields {
					namedType := vv.Type.NamedType
					if vv.Type.NamedType == "" {
						namedType = vv.Type.Elem.NamedType
						if vv.Type.Elem.NamedType == "" {
							namedType = vv.Type.Elem.Elem.NamedType
						}
					}
					switch gqlGenerate.Schema.Types[namedType].Kind {
					case "OBJECT":
						if !slices.Contains(skipObject, strings.Trim(vv.Name, " ")) {
							authRegexResult := authRegex.FindStringSubmatch(vv.Description)
							if len(authRegexResult) > 1 {
								if _, exists := authData.TypeNames[strings.ToUpper(namedType)]; !exists {
									authData.TypeNames[strings.ToUpper(namedType)] = namedType
								}
								if _, exists := authData.TypeNames[strings.ToUpper(v.Name)]; !exists {
									authData.TypeNames[strings.ToUpper(v.Name)] = v.Name
								}
								if _, exists := authData.Resolvers[strings.ToUpper(vv.Name)]; !exists {
									authData.Resolvers[strings.ToUpper(vv.Name)] = vv.Name
								}
								if _, exists := authData.Grants[strings.ToUpper(v.Name)]; !exists {
									authData.Grants[strings.ToUpper(v.Name)] = map[string]map[string][]string{}
								}
								if _, exists := authData.Grants[strings.ToUpper(v.Name)][strings.ToUpper(namedType)]; !exists {
									authData.Grants[strings.ToUpper(v.Name)][strings.ToUpper(namedType)] = map[string][]string{}
								}
								access := make([]string, 0)
								for _, v := range strings.Split(authRegexResult[1], ",") {
									access = append(access, strings.Trim(v, " "))
								}
								authData.Grants[strings.ToUpper(v.Name)][strings.ToUpper(namedType)][strings.ToUpper(vv.Name)] = access
							}
						}
					}
				}
			}
		}
	}
	generateTemplate(authData, gqlGenerate)
}

func generateTemplate(authData AuthData, gqlGenerate generate.GqlGenerate) {
	ut, err := template.New("enum.tmpl").Funcs(template.FuncMap{"StringsJoin": StringsJoin}).Parse(string(generate.Authtmpl))
	if err != nil {
		panic(err)
	}
	filePath := gqlGenerate.ModulePath + "/generate/" + gqlGenerate.LibPath + "/auth.go"
	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, 0770)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	var tmpl bytes.Buffer
	err = ut.Execute(&tmpl, authData)
	if err != nil {
		panic(err)
	}
	x, _ := format.Source(tmpl.Bytes())
	file.Write(x)

}
func StringsJoin(elem []string) (r string) {

	if elem == nil {
		r = "{}"
	} else {
		r = "{`" + strings.Join(elem, "`,`") + "`}"
	}
	return
}
