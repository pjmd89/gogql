package resolvers

import (
	"net/http"

	"github.com/fatih/structs"
	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type Schema *ast.Schema
type DataReturn interface{}
type Args map[string]interface{}
type Parent interface{}

type Directive interface {
	Invoke(args map[string]interface{}, typeName string, fieldName string) (DataReturn, definitionError.GQLError)
}
type DirectiveList map[string]interface{}

type Definition *ast.Definition
type RestInfo struct {
	Path      string
	PathSplit []string
	GET       map[string][]string
	POST      map[string][]string
	headers   map[string]string
	r         *http.Request
}
type Subscription struct {
	socketId       string
	resolverId     string
	subscriptionId int
}
type ResolverInfo struct {
	Operation         string
	Resolver          string
	Args              Args
	Parent            Parent
	Directives        DirectiveList
	TypeName          string
	ParentTypeName    *string
	SubscriptionValue interface{}
	SessionID         string
	RestInfo          *RestInfo
}
type ObjectTypeInterface interface {
	Resolver(ResolverInfo) (DataReturn, definitionError.GQLError)
	Subscribe(ResolverInfo) bool
}
type Scalar interface {
	Assess(resolved ScalarResolved) (r interface{}, err definitionError.GQLError)
	Set(value interface{}) (r interface{}, err definitionError.GQLError)
}
type ScalarResolved struct {
	Value        interface{}
	ResolverName string
	Resolved     *structs.Struct
}
