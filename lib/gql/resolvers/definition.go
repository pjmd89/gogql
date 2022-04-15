package resolvers

import (
	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type Schema *ast.Schema;
type DataReturn interface{}
type Args map[string]interface{}
type Parent interface{}

type Directive interface{
	Invoke(args map[string]interface{},typeName string, fieldName string) (DataReturn, definitionError.Error)
}
type DirectiveList map[string]interface{}

type Definition *ast.Definition;
type Subscription struct{
	socketId string;
    resolverId string;
    subscriptionId int;
}
type ResolverInfo struct{
	Operation string
	Resolver string
	Args Args
	Parent Parent
	Directives DirectiveList
	TypeName string
	ParentTypeName *string
	SubscriptionValue interface{}
}
type ObjectTypeInterface interface{
    Resolver(ResolverInfo) (DataReturn, definitionError.Error)
	Subscribe(ResolverInfo) (bool)
}
type Scalar interface{
	Assess(value interface{}) (r interface{}, err definitionError.Error)
}