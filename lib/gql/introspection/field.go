package introspection

import "github.com/vektah/gqlparser/v2/ast"
type Field struct{
	Name				string							`gql:"name=name"`
	Description			*string							`gql:"name=description"`
	Args				ast.ArgumentDefinitionList		`gql:"name=args"`
	Type 				*ast.Type						`gql:"name=type"`
	IsDeprecated		bool							`gql:"name=isDeprecated"`
	DeprecationReason 	*string							`gql:"name=deprecationReason"`

}
