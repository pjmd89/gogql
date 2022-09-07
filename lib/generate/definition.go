package generate

type AttrDef struct {
	Name    string
	Type    string
	IsArray bool
	BSONTag string
	GQLTag  string
}
type ModelDef struct {
	Name        string
	PackageName string
	Attr        []AttrDef
	IsUseID     bool
}
type EnumAttrDef struct {
	Name  string
	Value string
}
type EnumDef struct {
	Name        string
	PackageName string
	Attr        []EnumAttrDef
}
type UnionDef struct {
	Name        string
	PackageName string
	Types       []string
}
type ScalarDef struct {
	Name        string
	PackageName string
}
