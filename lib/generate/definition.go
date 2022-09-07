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
