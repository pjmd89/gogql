package introspection

import "github.com/vektah/gqlparser/v2/ast"

type Fields 		[]*Field
type Interfaces 	[]*Type
type PossibleTypes 	[]*Type
type EnumValues 	[]*EnumValue
type InputFields 	[]*InputValue	
type Type struct{
	Kind				TypeKind		`gql:"name=kind"`
	Name 				*string			`gql:"name=name"`
	Description 		*string			`gql:"name=description"`
	// OBJECT and INTERFACE only
	Fields 				*ast.FieldList	`gql:"name=fields"`
	// OBJECT only
	Interfaces	 		*[]string		`gql:"name=interfaces"`
	// INTERFACE and UNION only
	PossibleTypes 		*PossibleTypes	`gql:"name=possibleTypes"`
	// ENUM only
	EnumValues			*EnumValues		`gql:"name=enumValues"`
	// INPUT_OBJECT only
	InputFields			*InputFields	`gql:"name=inputFields"`
	// NON_NULL and LIST only
	OfType				*ast.Type		`gql:"name=ofType"`
}

func(o *Fields) IncludeDeprecated(includeDeprecate bool) *Fields{
	return o;
}

func(o *EnumValues) IncludeDeprecated(include bool) *EnumValues{
	return o;
}

