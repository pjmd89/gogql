package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gogql/lib/gql"
	"golang.org/x/exp/slices"
)

var (
	//go:embed templates/model.tmpl
	modeltmpl []byte
	//go:embed templates/enum.tmpl
	enumtmpl []byte
	//go:embed templates/union.tmpl
	uniontmpl []byte
	//go:embed templates/scalar.tmpl
	scalartmpl []byte
)
var omitObject = []string{
	"__Directive",
	"__EnumValue",
	"__Field",
	"__InputValue",
	"__Schema",
	"__Type",
	"__TypeKind",
	"__DirectiveLocation",
	"Query",
	"Mutation",
	"query",
	"mutation",
}
var typeToChange = map[string]string{
	"Int":     "int64",
	"Float":   "float64",
	"ID":      "primitive.ObjectID",
	"String":  "string",
	"Boolean": "bool",
}
var omitScalarTypes = []string{
	"Int",
	"Float",
	"ID",
	"String",
	"Boolean",
}
var indexIDName = []string{
	"_id",
	"id",
}

func main() {
	var scheme string
	var modelPath string
	flag.StringVar(&scheme, "scheme", "", "Ruta de la carpeta contenedora del esquema de GraphQL")
	flag.StringVar(&modelPath, "model-path", "", "Ruta donde se guardaran los modelos generados")
	flag.Parse()
	if scheme != "" && modelPath != "" {
		generateSchema(scheme, modelPath)
	}
}

func generateSchema(scheme string, modelPath string) {
	typeRegex := regexp.MustCompile(`-type ([^\n]+)`)
	pointerRegex := regexp.MustCompile(`-pointer[^\n]?`)
	defaultRegex := regexp.MustCompile(`-default ([^\n]+)`)
	createdRegex := regexp.MustCompile(`-created[^\n]?`)
	updatedRegex := regexp.MustCompile(`-updated[^\n]?`)
	mt, err := template.New("model.tmpl").Parse(string(modeltmpl))
	if err != nil {
		panic(err)
	}
	et, err := template.New("enum.tmpl").Parse(string(enumtmpl))
	if err != nil {
		panic(err)
	}
	ut, err := template.New("union.tmpl").Parse(string(uniontmpl))
	if err != nil {
		panic(err)
	}
	st, err := template.New("scalar.tmpl").Parse(string(scalartmpl))
	if err != nil {
		panic(err)
	}
	gql := gql.Init("", scheme)
	parent := filepath.Base(modelPath)
	values := gql.GetTypes()

	for key, value := range values {
		key = strings.Title(key)
		isID := false
		filename := strings.ToLower(key) + ".go"
		if !slices.Contains(omitObject, key) {
			switch value.Kind {
			case "OBJECT":
				structDef := &generate.ModelDef{}
				structDef.Attr = make([]generate.AttrDef, 0)
				structDef.Name = key
				structDef.PackageName = parent
				for _, vValue := range value.Fields {
					attrStruct := generate.AttrDef{}
					attrStruct.Name = strings.Title(vValue.Name)
					attrStruct.Type, isID = getNamedType(vValue.Type.NamedType)
					bsonTag := make([]string, 0)
					gqlTag := make([]string, 0)
					var namedTyped string
					vValueName := vValue.Name
					attrStruct.IsArray = false
					if slices.Contains(indexIDName, vValue.Name) {
						attrStruct.Name = "Id"
						vValueName = "_id"
					}
					bsonTag = append(bsonTag, vValueName)
					gqlTag = append(gqlTag, "name="+vValueName)
					if slices.Contains(indexIDName, vValue.Name) && vValue.Type.NamedType == "ID" {
						bsonTag = append(bsonTag, "omitempty")
						gqlTag = append(gqlTag, "id=true")
						structDef.IsUseID = true
					}

					if vValue.Type.Elem != nil {
						namedTyped, isID = getNamedType(vValue.Type.Elem.NamedType)
						attrStruct.Type = "[]" + namedTyped
						attrStruct.IsArray = true
					}
					if vValue.Type.Elem != nil && vValue.Type.Elem.NonNull == false {
						namedTyped, isID = getNamedType(vValue.Type.Elem.NamedType)
						attrStruct.Type = "[]" + namedTyped
						attrStruct.IsArray = true
					}
					if vValue.Type.Elem != nil && vValue.Type.Elem.NonNull == true && vValue.Type.Elem.Elem != nil {
						namedTyped, isID = getNamedType(vValue.Type.Elem.Elem.NamedType)
						attrStruct.Type = "[]" + namedTyped
						attrStruct.IsArray = true
					}
					typeRegexResult := typeRegex.FindStringSubmatch(vValue.Description)
					if len(typeRegexResult) > 1 {
						namedTyped, isID = getNamedType(typeRegexResult[1])
						if attrStruct.IsArray {

							attrStruct.Type = "[]" + namedTyped
						} else {
							attrStruct.Type = namedTyped
						}
					}
					if isID {
						structDef.IsUseID = true
					}
					pointerRegexResult := pointerRegex.MatchString(vValue.Description)
					if pointerRegexResult {
						attrStruct.Type = "*" + attrStruct.Type
					}
					if structDef.IsUseID {
						gqlTag = append(gqlTag, "objectID=true")
						structDef.IsUseID = true
					}
					defaultRegexResult := defaultRegex.FindStringSubmatch(vValue.Description)
					if len(defaultRegexResult) > 1 {
						gqlTag = append(gqlTag, "default="+defaultRegexResult[1])
					}
					createdRegexResult := createdRegex.MatchString(vValue.Description)
					if createdRegexResult {
						gqlTag = append(gqlTag, "created=true")
					}
					updatedRegexResult := updatedRegex.MatchString(vValue.Description)
					if updatedRegexResult {
						gqlTag = append(gqlTag, "updated=true")
					}
					attrStruct.BSONTag = strings.Join(bsonTag, ",")
					attrStruct.GQLTag = strings.Join(gqlTag, ",")
					structDef.Attr = append(structDef.Attr, attrStruct)
				}
				dir := filepath.Dir(modelPath + "/model_" + filename)
				os.MkdirAll(dir, 0770)
				modelFile, err := os.Create(modelPath + "/model_" + filename)
				if err != nil {
					panic(err)
				}
				if structDef.Name == "State" {
					fmt.Println("")
				}
				var tmpl bytes.Buffer
				err = mt.Execute(&tmpl, structDef)
				if err != nil {
					panic(err)
				}
				x, _ := format.Source(tmpl.Bytes())
				modelFile.Write(x)
				break
			case "ENUM":
				structDef := &generate.EnumDef{}
				structDef.Attr = make([]generate.EnumAttrDef, 0)
				structDef.Name = key
				structDef.PackageName = parent
				for _, vValue := range value.EnumValues {
					enumAttrDef := generate.EnumAttrDef{}
					enumAttrDef.Name = strings.ToUpper(key + "_" + vValue.Name)
					enumAttrDef.Value = vValue.Name
					structDef.Attr = append(structDef.Attr, enumAttrDef)
				}
				dir := filepath.Dir(modelPath + "/enum_" + filename)
				os.MkdirAll(dir, 0770)
				modelFile, err := os.Create(modelPath + "/enum_" + filename)
				if err != nil {
					panic(err)
				}
				var tmpl bytes.Buffer
				err = et.Execute(&tmpl, structDef)
				if err != nil {
					panic(err)
				}
				x, _ := format.Source(tmpl.Bytes())
				modelFile.Write(x)
				if err != nil {
					panic(err)
				}
				break
			case "SCALAR":
				if !slices.Contains(omitScalarTypes, key) {
					scalarDef := &generate.ScalarDef{}
					scalarDef.Name = key
					scalarDef.PackageName = parent
					dir := filepath.Dir(modelPath + "/scalar_" + filename)
					os.MkdirAll(dir, 0770)
					modelFile, err := os.Create(modelPath + "/scalar_" + filename)
					if err != nil {
						panic(err)
					}
					var tmpl bytes.Buffer
					err = st.Execute(&tmpl, scalarDef)
					if err != nil {
						panic(err)
					}
					x, _ := format.Source(tmpl.Bytes())
					modelFile.Write(x)
					if err != nil {
						panic(err)
					}
				}
				break
			case "UNION":
				unionDef := &generate.UnionDef{}
				unionDef.Name = key
				unionDef.PackageName = parent
				for _, vValue := range value.Types {
					unionDef.Types = append(unionDef.Types, strings.Title(vValue))
				}
				dir := filepath.Dir(modelPath + "/union_" + filename)
				os.MkdirAll(dir, 0770)
				modelFile, err := os.Create(modelPath + "/union_" + filename)
				if err != nil {
					panic(err)
				}
				var tmpl bytes.Buffer
				err = ut.Execute(&tmpl, unionDef)
				if err != nil {
					panic(err)
				}
				x, _ := format.Source(tmpl.Bytes())
				modelFile.Write(x)
				if err != nil {
					panic(err)
				}
				break
			}
		}
	}

}
func getNamedType(namedType string) (r string, isID bool) {
	if _, ok := typeToChange[namedType]; ok {
		r = typeToChange[namedType]
		if namedType == "ID" {
			isID = true
		}
	} else {
		r = strings.Title(namedType)
	}
	return
}

//go:generate go run main.go -scheme=$SCHEME -model-path=$MODELPATH
