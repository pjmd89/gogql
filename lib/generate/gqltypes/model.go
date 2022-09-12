package gqltypes

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gqlparser/v2/ast"
	"golang.org/x/exp/slices"
)

func NewModel(render generate.GqlGenerate, key string, value *ast.Definition, types map[string]*ast.Definition) generate.ModelDef {
	typeRegex := regexp.MustCompile(`-type ([^\n]+)`)
	gqlOmitRegex := regexp.MustCompile(`-gqlOmit[^\n]?`)
	bsonOmitRegex := regexp.MustCompile(`-bsonOmit[^\n]?`)
	pointerRegex := regexp.MustCompile(`-pointer[^\n]?`)
	defaultRegex := regexp.MustCompile(`-default ([^\n]+)`)
	nestedRegex := regexp.MustCompile(`-nested[^\n]?`)
	createdRegex := regexp.MustCompile(`-created[^\n]?`)
	updatedRegex := regexp.MustCompile(`-updated[^\n]?`)
	structDef := generate.ModelDef{}
	structDef.Attr = make([]generate.AttrDef, 0)
	structDef.Name = strings.Title(key)
	structDef.Name = strings.Title(key)
	structDef.RealName = key
	structDef.PackageName = render.ModelPath
	bsonOmit := false
	gqlOmit := false
	isID := false
	unionCount := 0
	unionList := []string{}
	for _, vValue := range value.Fields {
		attrStruct := generate.AttrDef{}
		namedType := vValue.Type.NamedType
		if namedType == "" {
			namedType = vValue.Type.Elem.NamedType
		}
		attrStruct.Name = strings.Title(vValue.Name)
		attrStruct.Type, isID = generate.GetNamedType(namedType)
		bsonTag := make([]string, 0)
		gqlTag := make([]string, 0)
		unionInstance := ""
		if types[namedType] != nil {
			switch types[namedType].Kind {
			case "ENUM":
				structDef.IsUseEnum = true
				attrStruct.Type = render.EnumPath + "." + attrStruct.Type
				break
			case "UNION":
				structDef.IsUseUnion = true
				unionCount++
				unionInstance = attrStruct.Type + strconv.Itoa(unionCount)
				unionList = append(unionList, unionInstance+" "+attrStruct.Type)
				attrStruct.Type = unionInstance
				break
			case "SCALAR":
				if !slices.Contains(generate.OmitScalarTypes, namedType) {
					structDef.IsUseScalar = true
					attrStruct.Type = render.ScalarPath + "." + attrStruct.Type
				}
				break
			}
		}

		var namedTyped string
		vValueName := vValue.Name
		attrStruct.IsArray = false
		if slices.Contains(generate.IndexIDName, vValue.Name) {
			attrStruct.Name = "Id"
			vValueName = "_id"
		}
		bsonTag = append(bsonTag, vValueName)
		gqlTag = append(gqlTag, "name="+vValueName)
		if slices.Contains(generate.IndexIDName, vValue.Name) && namedType == "ID" {
			bsonTag = append(bsonTag, "omitempty")
			gqlTag = append(gqlTag, "id=true")
			structDef.IsUseID = true
		}

		if vValue.Type.Elem != nil {
			namedTyped, isID = generate.GetNamedType(vValue.Type.Elem.NamedType)
			if unionInstance != "" {
				//namedTyped = unionInstance
			}
			attrStruct.Type = "[]" + namedTyped
			attrStruct.IsArray = true
		}
		if vValue.Type.Elem != nil && vValue.Type.Elem.NonNull == false {
			namedTyped, isID = generate.GetNamedType(vValue.Type.Elem.NamedType)
			if unionInstance != "" {
				//namedTyped = unionInstance
			}
			attrStruct.Type = "[]" + namedTyped
			attrStruct.IsArray = true
		}
		if vValue.Type.Elem != nil && vValue.Type.Elem.NonNull == true && vValue.Type.Elem.Elem != nil {
			namedTyped, isID = generate.GetNamedType(vValue.Type.Elem.Elem.NamedType)
			if unionInstance != "" {
				//namedTyped = unionInstance
			}
			attrStruct.Type = "[]" + namedTyped
			attrStruct.IsArray = true
		}
		typeRegexResult := typeRegex.FindStringSubmatch(vValue.Description)
		if len(typeRegexResult) > 1 {
			namedTyped, isID = generate.GetNamedType(typeRegexResult[1])
			if unionInstance != "" {
				//namedTyped = unionInstance
			}
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
		if isID {
			gqlTag = append(gqlTag, "objectID=true")
			structDef.IsUseID = true
		}
		if gqlOmitRegex.MatchString(vValue.Description) {
			gqlOmit = true
		}
		if bsonOmitRegex.MatchString(vValue.Description) {
			bsonOmit = true
		}
		nestedResult := nestedRegex.MatchString(vValue.Description)
		if nestedResult {
			gqlTag = append(gqlTag, "nested=true")
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
		if bsonOmit {
			attrStruct.BSONTag = "-"
		}
		attrStruct.GQLTag = strings.Join(gqlTag, ",")
		if gqlOmit {
			attrStruct.GQLTag = "omit=true"
		}

		structDef.Attr = append(structDef.Attr, attrStruct)
	}
	if unionCount > 0 {
		//structDef.Name = structDef.Name + "[" + strings.Join(unionList, ",") + "]"
	}
	structDef.FilePath = render.ModulePath + "/generate/" + render.ModelPath + "/" + strings.ToLower(key) + ".go"
	structDef.ModelPath = render.ModuleName + "/" + render.ModelPath
	structDef.ScalarPath = render.ModuleName + "/" + render.ResolverPath + "/" + render.ScalarPath
	structDef.UnionPath = render.ModuleName + "/" + render.ModelPath + "/" + render.UnionPath
	structDef.EnumPath = render.ModuleName + "/" + render.ModelPath + "/" + render.EnumPath
	objectTypeBase := render.ModulePath + "/generate/" + render.ResolverPath + "/" + render.ObjecttypePath + "/" + strings.ToLower(key)
	structDef.ObjectTypePath = objectTypeBase + "/definition.go"
	structDef.QueryPath = objectTypeBase + "/queries.go"
	structDef.MutationPath = objectTypeBase + "/mutations.go"
	structDef.SubscriptionPath = objectTypeBase + "/subscriptions.go"
	return structDef
}
func ModelTmpl(types generate.RenderTypes) {
	mt, err := template.New("model.tmpl").Parse(string(generate.Modeltmpl))
	if err != nil {
		panic(err)
	}
	for _, v := range types.ModelType {
		dir := filepath.Dir(v.FilePath)
		os.MkdirAll(dir, 0770)
		modelFile, err := os.Create(v.FilePath)
		if err != nil {
			panic(err)
		}

		var tmpl bytes.Buffer
		err = mt.Execute(&tmpl, v)
		if err != nil {
			panic(err)
		}
		x, err := format.Source(tmpl.Bytes())
		if err != nil {
			log.Fatal(err.Error())
		}
		modelFile.Write(x)
	}
}
