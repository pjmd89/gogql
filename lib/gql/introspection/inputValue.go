package introspection

import "github.com/vektah/gqlparser/v2/ast"
type Args []*InputValue
type InputValue struct{
	Name				string		`gql:"name=name,storage=name"`
	Description			*string		`gql:"name=description,storage=description"`
	Type				*ast.Type	`gql:"name=type,storage=type"`
	DefaultValue		*string		`gql:"namez=defaultValue,storage=defaultValue"`
	ParentType			string		`gql:"storage=parentType"`
	FindIn				string		`gql:"storage=findIn"`
}