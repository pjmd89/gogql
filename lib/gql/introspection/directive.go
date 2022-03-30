package introspection

import "github.com/pjmd89/gqlparser/v2/ast"

type Locations []*DirectiveLocation

type Directive struct{
	Name 			string						`gql:"name=name"`
	Description 	*string						`gql:"name=description"`
	Locations 		[]ast.DirectiveLocation		`gql:"name=locations"`
	Args			ast.ArgumentDefinitionList	`gql:"name=args"`
	ParentType		string						`gql:"storage=parentType"`
	Directive		bool						`gql:"storage=directive"`
}