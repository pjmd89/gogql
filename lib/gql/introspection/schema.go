package introspection

import "github.com/pjmd89/gqlparser/v2/ast"

type Directives []*Directive

type Schema struct{
	Types 				map[string]*ast.Definition			`gql:"name=types"`
	QueryType			*ast.Definition						`gql:"name=queryType"`
	MutationType		*ast.Definition						`gql:"name=mutationType"`
	SubscriptionType 	*ast.Definition						`gql:"name=subscriptionType"`
	Directives			map[string]*ast.DirectiveDefinition	`gql:"name=directives"`
}
func(o *Schema) SetDef(){

}