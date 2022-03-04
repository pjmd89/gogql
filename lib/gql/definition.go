package gql

import (
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gogql/lib/http"
	"github.com/vektah/gqlparser/v2/ast"
)

type ObjectTypes map[string]resolvers.ObjectTypeInterface;
type Directives map[string]resolvers.Directive;
type variables map[string]interface{};
type gql struct{
	serverName 		string
	schema 			*ast.Schema
    objectTypes 	ObjectTypes
    variables       variables
    Directives      Directives
    Session         *http.Session
}
type HttpRequest struct{
    Query         	string                 `json:"query"`
    Variables     	map[string]interface{} `json:"variables,omitempty"`
    OperationName 	string                 `json:"operationName,omitempty"`
}

type HttpResponse struct{
    Data    		string `json:"data"`
    //Errors  		[]*lib.GqlError `json:"errors,omitempty"`
}
