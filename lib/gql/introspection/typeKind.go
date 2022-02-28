package introspection

type TypeKind int

const (
	TYPEKIND_ERROR 					TypeKind = iota
	TYPEKIND_SCALAR
	TYPEKIND_OBJECT
	TYPEKIND_INTERFACE
	TYPEKIND_UNION
	TYPEKIND_ENUM
	TYPEKIND_INPUT_OBJECT
	TYPEKIND_LIST
	TYPEKIND_NON_NULL
)
func SetTypeKind(typekind string) *TypeKind {
	var tk TypeKind;
	
	switch typekind {
		case "SCALAR":
			tk = TYPEKIND_SCALAR;
			break;
		case "OBJECT":
			tk = TYPEKIND_OBJECT;
			break;
		case "INTERFACE":
			tk = TYPEKIND_INTERFACE;
			break;
		case "UNION":
			tk = TYPEKIND_UNION;
			break;
		case "ENUM":
			tk = TYPEKIND_ENUM;
			break;
		case "INPUT_OBJECT":
			tk = TYPEKIND_INPUT_OBJECT;
			break;
		case "LIST":
			tk = TYPEKIND_LIST;
			break;
		case "NON_NULL":
			tk = TYPEKIND_NON_NULL;
			break;
		default:
			tk = TYPEKIND_ERROR;
	}
	return &tk;
}

func(o TypeKind) String() string{
	return [...]string{
		"ERROR",
		"SCALAR",
		"OBJECT",
		"INTERFACE",
		"UNION",
		"ENUM",
		"INPUT_OBJECT",
		"LIST",
		"NON_NULL",
	}[o];
}
func(o TypeKind) Is(value string) bool {
	return o.String() == value;
}
func (o TypeKind) EnumIndex() int {
	return int(o)
}
func (o TypeKind) EnumText( x TypeKind) string {
	return x.String()
}