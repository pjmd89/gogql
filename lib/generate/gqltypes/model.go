package gqltypes

import (
	"bytes"
	"fmt"
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

func NewModel(render generate.GqlGenerate, key string, value *ast.Definition, types map[string]*ast.Definition, driverDB DriverDB) generate.ModelDef {
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
	isID := false
	unionCount := 0
	unionList := []string{}
	structDef.GQLFile = strings.Replace(value.Position.Src.Name, "../../", "../", 1) + ":" + strconv.Itoa(value.Position.Line)
	for _, vValue := range value.Fields {
		bsonOmit := false
		gqlOmit := false
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
				//unionInstance = attrStruct.Type + strconv.Itoa(unionCount)
				unionInstance = attrStruct.Type
				unionList = append(unionList, unionInstance+" "+attrStruct.Type)
				attrStruct.Type = unionInstance
				break
			case "SCALAR":
				if !slices.Contains(generate.OmitScalarTypes, namedType) {
					if attrStruct.Type != "ID" {
						structDef.IsUseScalar = true
					} else if driverDB == DRIVERDB_NONE && attrStruct.Type == "ID" {
						structDef.IsUseScalar = true
					}
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
			attrStruct.Type, structDef.DriverDB, structDef.IsDriverDB = renderID(render.ScalarPath, "ID", driverDB)
			structDef.IsUseID = true
		}

		if vValue.Type.Elem != nil {
			namedTyped, isID = generate.GetNamedType(vValue.Type.Elem.NamedType)
			if unionInstance != "" {
				//namedTyped = unionInstance
			}
			attrStruct.Type, structDef.DriverDB, structDef.IsDriverDB = renderID(render.ScalarPath, namedTyped, driverDB)
			attrStruct.Type = "[]" + attrStruct.Type
			attrStruct.IsArray = true
		}
		if vValue.Type.Elem != nil && vValue.Type.Elem.NonNull == false {
			namedTyped, isID = generate.GetNamedType(vValue.Type.Elem.NamedType)
			if unionInstance != "" {
				//namedTyped = unionInstance
			}
			attrStruct.Type, structDef.DriverDB, structDef.IsDriverDB = renderID(render.ScalarPath, namedTyped, driverDB)
			attrStruct.Type = "[]" + attrStruct.Type
			attrStruct.IsArray = true
		}
		if vValue.Type.Elem != nil && vValue.Type.Elem.NonNull == true && vValue.Type.Elem.Elem != nil {
			namedTyped, isID = generate.GetNamedType(vValue.Type.Elem.Elem.NamedType)
			if unionInstance != "" {
				//namedTyped = unionInstance
			}
			attrStruct.Type, structDef.DriverDB, structDef.IsDriverDB = renderID(render.ScalarPath, namedTyped, driverDB)
			attrStruct.Type = "[]" + attrStruct.Type
			attrStruct.IsArray = true
		}
		typeRegexResult := typeRegex.FindStringSubmatch(vValue.Description)
		if len(typeRegexResult) > 1 {
			xtypeRegex := regexp.MustCompile(`\[([^\]]+)\]`)
			xtypeRegexResult := xtypeRegex.FindStringSubmatch(typeRegexResult[1])
			newNamedType := typeRegexResult[1]
			if len(xtypeRegexResult) > 1 {
				newNamedType = xtypeRegexResult[1]
				attrStruct.IsArray = true
			}
			namedTyped, isID = generate.GetNamedType(newNamedType)
			if unionInstance != "" {
				//namedTyped = unionInstance
			}
			if attrStruct.IsArray {
				attrStruct.Type, structDef.DriverDB, structDef.IsDriverDB = renderID(render.ScalarPath, namedTyped, driverDB)
				attrStruct.Type = "[]" + attrStruct.Type
			} else {
				attrStruct.Type, structDef.DriverDB, structDef.IsDriverDB = renderID(render.ScalarPath, namedTyped, driverDB)
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
func renderID(scalarPath string, attrName string, driverDB DriverDB) (r string, dp string, isDB bool) {

	switch driverDB {
	case DRIVERDB_NONE:
		if attrName == "ID" {
			r = scalarPath + "." + attrName
		}
	case DRIVERDB_MONGO:
		r = "primitive.ObjectID"
		dp = "go.mongodb.org/mongo-driver/bson/primitive"
		isDB = true
	default:
		r = attrName
	}

	return
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
			fmt.Printf("Error creando Model %s: %s - ", v.FilePath, err)
			panic(err)
		}

		var tmpl bytes.Buffer
		err = mt.Execute(&tmpl, v)
		if err != nil {
			fmt.Printf("Error creando model %s: %s - ", v.Name, err)
			panic(err)
		}
		x, err := format.Source(tmpl.Bytes())
		if err != nil {
			fmt.Printf("Error creando model %s: %s - ", v.Name, err)
			log.Fatal(err.Error())
		}
		modelFile.Write(x)
	}
}
