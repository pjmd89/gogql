package gql

import (
	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type ObjectTypes map[string]resolvers.ObjectTypeInterface
type Directives map[string]resolvers.Directive
type Scalars map[string]resolvers.Scalar
type OperationID string
type SocketID string
type EventID string
type RequestID string
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
type gql struct {
	serverName       string
	schema           *ast.Schema
	objectTypes      ObjectTypes
	directives       Directives
	scalars          Scalars
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
	Data   string                     `json:"data,omitempty"`
	Errors []definitionError.GQLError `json:"errors,omitempty"`
}

type DefaultArguments struct {
	Name    string
	IsArray bool
	Value   interface{}
	NonNull bool
	Kind    string
	Type    string
}
