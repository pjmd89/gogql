package resolvers

import (
	"github.com/pjmd89/gogql/lib/gql/pubsub"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type Schema *ast.Schema;
type DataReturn interface{}
type Args map[string]interface{}
type Parent interface{}
type Subscription *pubsub.Subscription
type Directive interface{
	Invoke(args map[string]interface{},typeName string, fieldName string) DataReturn
}
type DirectiveList map[string]interface{}

type Definition *ast.Definition;
type ResolverInfo struct{
	Operation string
	Resolver string
	Args Args
	Parent Parent
	Directives DirectiveList
	TypeName string
	ParentTypeName *string
	Subscription Subscription
}
type ObjectTypeInterface interface{
    Resolver(ResolverInfo) DataReturn
	Subscribe(ResolverInfo) (bool, Subscription)
}
type Scalar interface{
	Assess(value interface{}) (r interface{}, err error)
}