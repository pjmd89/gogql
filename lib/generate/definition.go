package generate

import (
	_ "embed"
	"strings"

	"github.com/pjmd89/gqlparser/v2/ast"
)

type GqlGenerate struct {
	Schema         *ast.Schema
	ModuleName     string
	ModulePath     string
	ModelPath      string
	ResolverPath   string
	UnionPath      string
	ScalarPath     string
	EnumPath       string
	ObjecttypePath string
	LibPath        string
	SchemaPath     string
}
type RenderTypes struct {
	ModelType     []ModelDef
	ObjectType    []ObjectTypeDef
	EnumType      []EnumDef
	UnionType     UnionDef
	ScalarType    []ScalarDef
	IsScalar      bool
	ScalarPath    string
	MainPath      string
	SchemaPath    string
	LibConfigPath string
	ConfigDB      string
	ConfigHTTP    string
	ConfigJSON    string
	ModuleName    string
}
type ObjectTypeDef struct {
	Name                  string
	PackageName           string
	RealName              string
	ModelPath             string
	ModuleName            string
	DefinitionPath        string
	FilePath              string
	QueryPath             string
	MutationPath          string
	SubscriptionPath      string
	HasQueries            bool
	HasMutations          bool
	HasSubscriptions      bool
	QueryResolvers        []map[string]string
	MutationResolvers     []map[string]string
	SubscriptionResolvers []map[string]string
}
type ModelDef struct {
	Name             string
	RealName         string
	PackageName      string
	Attr             []AttrDef
	DriverDB         string
	IsDriverDB       bool
	IsUseID          bool
	IsUseUnion       bool
	IsUseScalar      bool
	IsUseEnum        bool
	FilePath         string
	ModelPath        string
	ScalarPath       string
	UnionPath        string
	EnumPath         string
	ObjectTypePath   string
	QueryPath        string
	MutationPath     string
	SubscriptionPath string
	GQLFile          string
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
	VarName     string
	PackageName string
	TypeName    string
	FilePath    string
}

var (
	//go:embed templates/auth.tmpl
	Authtmpl []byte
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
	//go:embed templates/queries.tmpl
	Queriestmpl []byte
	//go:embed templates/mutations.tmpl
	Mutationstmpl []byte
	//go:embed templates/subscriptions.tmpl
	Subscriptionstmpl []byte
	//go:embed templates/main.tmpl
	Maintmpl []byte
	//go:embed templates/libconfig.tmpl
	LibConfigtmpl []byte
	//go:embed templates/configjson.tmpl
	ConfigJSONtmpl []byte
	//go:embed templates/configdb.tmpl
	ConfigDBtmpl []byte
	//go:embed templates/confighttp.tmpl
	ConfigHTTPtmpl []byte
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
	"String":  "string",
	"Boolean": "bool",
}

func GetNamedType(namedType string) (r string, isID bool) {
	if _, ok := typeToChange[namedType]; ok {
		r = typeToChange[namedType]
	} else {
		if namedType == "ID" {
			isID = true
		}
		r = strings.Title(namedType)
	}
	return
}
