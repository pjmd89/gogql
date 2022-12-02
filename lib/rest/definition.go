package rest

import "github.com/pjmd89/gogql/lib/resolvers"

type ObjectTypes map[string]ObjectType
type ObjectType struct {
	Alias      string
	ObjectType resolvers.ObjectTypeInterface
}
type rest struct {
	serverName  string
	objectTypes ObjectTypes
}
