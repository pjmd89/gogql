package generate

import (
	_ "embed"
	"strings"
)

type GqlGenerate struct {
	SchemaPath     string
	ModuleName     string
	ModulePath     string
	ModelPath      string
	ResolverPath   string
	UnionPath      string
	ScalarPath     string
	EnumPath       string
	ObjecttypePath string
}
type RenderTypes struct {
	ModelType  []ModelDef
	ObjectType map[string]ObjectTypeDef
	EnumType   []EnumDef
	UnionType  UnionDef
	ScalarType []ScalarDef
}
type ObjectTypeDef struct{}
type ModelDef struct {
	Name        string
	PackageName string
	Attr        []AttrDef
	IsUseID     bool
	IsUseUnion  bool
	IsUseScalar bool
	IsUseEnum   bool
	FilePath    string
	ScalarPath  string
	UnionPath   string
	EnumPath    string
}

type AttrDef struct {
	Name    string
	Type    string
	IsArray bool
	BSONTag string
	GQLTag  string
}
type EnumAttrDef struct {
	Name  string
	Value string
}
type EnumDef struct {
	Name        string
	PackageName string
	Attr        []EnumAttrDef
	FilePath    string
}
type UnionDef struct {
	PackageName string
	Type        []UnionAttrDef
	FilePath    string
}
type UnionAttrDef struct {
	Name  string
	Types string
}
type ScalarDef struct {
	Name        string
	PackageName string
	TypeName    string
	FilePath    string
}

var (
	//go:embed templates/model.tmpl
	Modeltmpl []byte
	//go:embed templates/enum.tmpl
	Enumtmpl []byte
	//go:embed templates/union.tmpl
	Uniontmpl []byte
	//go:embed templates/scalar.tmpl
	Scalartmpl []byte
	//go:embed templates/objecttype_definition.tmpl
	Objecttypetmpl []byte
)
var OmitObject = []string{
	"__Directive",
	"__EnumValue",
	"__Field",
	"__InputValue",
	"__Schema",
	"__Type",
	"__TypeKind",
	"__DirectiveLocation",
}
var OmitScalarTypes = []string{
	"Int",
	"Float",
	"ID",
	"String",
	"Boolean",
}
var IndexIDName = []string{
	"_id",
	"id",
}
var typeToChange = map[string]string{
	"Int":     "int64",
	"Float":   "float64",
	"ID":      "primitive.ObjectID",
	"String":  "string",
	"Boolean": "bool",
}

func GetNamedType(namedType string) (r string, isID bool) {
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
