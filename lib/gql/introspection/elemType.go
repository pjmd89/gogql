package introspection

type ElemType struct{
	NonNull 	bool		`gql:"storage=nonNull"`
	NamedType 	string		`gql:"storage=namedType"`
	Elem 		interface{}	`gql:"storage=elem"`
}