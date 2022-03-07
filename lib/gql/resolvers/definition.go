package resolvers

import (
	"github.com/vektah/gqlparser/v2/ast"
)

type Schema *ast.Schema;
type DataReturn interface{}
type Args map[string]interface{}
type Parent interface{}
type Blank struct{}

type Directive interface{
	Invoke(args map[string]interface{},typeName string, fieldName string, directiveDefinition *ast.DirectiveDefinition) DataReturn
}
type DirectiveList map[string]interface{}

type Definition *ast.Definition;
type Storage interface{}
type ObjectTypeInterface interface{
    Resolver( string, Args, Parent, DirectiveList, string) DataReturn
}