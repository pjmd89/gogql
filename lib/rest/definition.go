package rest

import "github.com/pjmd89/gogql/lib/resolvers"

type ObjectTypes map[string]ObjectType
type ObjectType struct {
	Alias      string
	ObjectType resolvers.ObjectTypeInterface
}
type rest struct {
	objectTypes ObjectTypes
}

type WebSocketRequest struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func Init() (r *rest) {
	r = &rest{}
	r.objectTypes = make(map[string]ObjectType)
	return
}
