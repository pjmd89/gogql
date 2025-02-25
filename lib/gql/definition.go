package gql

import (
	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/resolvers"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type ObjectTypes map[string]resolvers.ObjectTypeInterface
type Directives map[string]resolvers.Directive
type Scalars map[string]resolvers.Scalar
type OperationID string
type SocketID string
type EventID string
type RequestID string
type ResolverName string
type TypeName string
type Access []string
type Grant map[TypeName]map[TypeName]map[ResolverName]Access

type AuthorizateInfo struct {
	Operation string
	SrcType   TypeName
	DstType   TypeName
	Resolver  ResolverName
	SessionID string
}

type Subscription struct {
	channel     chan bool
	eventID     EventID
	operationID OperationID
	socketID    SocketID
	requestID   RequestID
	messageType int
	value       interface{}
}
type SubscriptionClose struct{}
type SourceEvents struct {
	subscriptionEvents map[OperationID]chan interface{}
	operationEvents    map[OperationID]map[EventID]*Subscription
}
type Gql struct {
	serverName       string
	schema           *ast.Schema
	objectTypes      ObjectTypes
	directives       Directives
	scalars          Scalars
	OnIntrospection  func() (err definitionError.GQLError)
	OnAuthorizate    func(authInfo AuthorizateInfo) definitionError.GQLError
	OnScalarArgument func(scalarType string, value interface{}) (r interface{})
}
type HttpRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}
type WebSocketRequest struct {
	Id      string      `json:"id"`
	Type    string      `json:"type"`
	Payload HttpRequest `json:"payload"`
}

type HttpResponse struct {
	Data   string `json:"data,omitempty"`
	Errors string `json:"errors,omitempty"`
}

type DefaultArguments struct {
	Name    string
	IsArray bool
	Value   interface{}
	NonNull bool
	Kind    string
	Type    string
}
