package gql

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	gqlHttp "github.com/pjmd89/gogql/lib/http"
	"github.com/pjmd89/gogql/lib/resolvers"
	"github.com/pjmd89/gogql/lib/resolvers/directives"
	"github.com/pjmd89/gogql/lib/resolvers/objectTypes"
	"github.com/pjmd89/gogql/lib/resolvers/scalars"
	"github.com/pjmd89/goutils/systemutils"

	"github.com/pjmd89/gqlparser/v2"
	"github.com/pjmd89/gqlparser/v2/ast"
)

func Init(filesystem systemutils.FSInterface, folder string) *gql {
	gql := &gql{}
	gql.loadSchema(filesystem, folder)
	gql.objectTypes = make(map[string]resolvers.ObjectTypeInterface)
	gql.directives = make(map[string]resolvers.Directive)
	gql.scalars = make(map[string]resolvers.Scalar)

	gql.objectTypes["__Schema"] = objectTypes.NewSchema(gql.schema, gql.directives)
	gql.objectTypes["__Type"] = objectTypes.NewType(gql.schema, gql.directives)
	gql.objectTypes["__Field"] = objectTypes.NewField(gql.schema, gql.directives)
	gql.objectTypes["__EnumValue"] = objectTypes.NewEnumValue(gql.schema, gql.directives)
	gql.objectTypes["__InputValue"] = objectTypes.NewInputValue(gql.schema, gql.directives)
	gql.objectTypes["__Directive"] = objectTypes.NewDirective(gql.schema)
	/*
		gql.directives["include"] = directives.NewInclude(gql.schema);
		gql.directives["skip"] = directives.NewSkip(gql.schema);
	   //*/
	gql.directives["deprecated"] = directives.NewDeprecated(gql.schema)
	gql.scalars["ID"] = scalars.NewIDScalar()
	gql.scalars["Boolean"] = scalars.NewBoolScalar()
	gql.scalars["String"] = scalars.NewStringScalar()
	gql.scalars["Int"] = scalars.NewIntScalar()
	gql.scalars["Float"] = scalars.NewFloatScalar()
	gql.OnAuthenticate = OnAuthenticate
	gql.OnIntrospection = OnIntrospection
	return gql
}
func OnIntrospection() (err definitionError.GQLError) {
	return
}
func OnAuthenticate(operation string, srcType, dstType TypeName, resolver string) (err definitionError.GQLError) {
	return
}
func (o *gql) GetScalars() Scalars {
	return o.scalars
}
func (o *gql) GetSchema() *ast.Schema {
	return o.schema
}
func (o *gql) ObjectType(resolver string, object resolvers.ObjectTypeInterface) {
	o.objectTypes[resolver] = object
}
func (o *gql) Directive(resolver string, object resolvers.Directive) {
	o.directives[resolver] = object
}
func (o *gql) Scalar(resolver string, object resolvers.Scalar) {
	o.scalars[resolver] = object
}
func (o *gql) scanSchema(filesystem systemutils.FSInterface, folder string) (r []string) {
	files, _ := filesystem.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			r = append(r, o.scanSchema(filesystem, folder+"/"+file.Name())...)
		} else {
			r = append(r, folder+"/"+file.Name())
		}
	}
	return
}
func (o *gql) loadSchema(filesystem systemutils.FSInterface, folder string) {
	var schema []*ast.Source

	files := o.scanSchema(filesystem, folder)

	for _, file := range files {
		content, err := filesystem.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		text := string(content)
		schema = append(schema, &ast.Source{Name: file, Input: text, BuiltIn: true})
	}
	parser, err := gqlparser.LoadSchema(schema...)

	if err != nil {
		log.Fatal(err.Error())

	}
	if parser.Subscription != nil {
		for _, subs := range parser.Subscription.Fields {
			PubSub.createSubscriptionEvent(OperationID(subs.Name))
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	o.schema = parser
}
func (o *gql) GQLRenderSubscription(mt int, message []byte, socketId, sessionID string) {
	var request WebSocketRequest
	json.Unmarshal(message, &request)
	//var response *HttpResponse = &HttpResponse{};
	switch request.Type {
	case "connection_init":
		r := `{"type":"connection_ack","payload":{}}`
		gqlHttp.WriteWebsocketMessage(mt, socketId, []byte(r))
	case "subscribe":
		o.WebsocketResponse(request.Payload, socketId, RequestID(request.Id), mt, sessionID)
		//r = `{"id":"`+request.Id+`","type":"next","payload":`+response.Data+`}`;
	case "ping":
		r := `{"type":"pong","payload":{}}`
		gqlHttp.WriteWebsocketMessage(mt, socketId, []byte(r))
	case "complete":
		r := `{"id":"` + request.Id + `","type":"complete"}`
		gqlHttp.WriteWebsocketMessage(mt, socketId, []byte(r))
	default:

	}
}
func (o *gql) GQLRender(w http.ResponseWriter, r *http.Request, sessionID string) {
	var request HttpRequest
	json.NewDecoder(r.Body).Decode(&request)
	response := o.response(request, sessionID)
	str := "{\"data\":" + response.Data
	if response.Errors != "" {
		str = str + ",\"errors\":" + response.Errors
	}
	str = str + "}"
	fmt.Fprint(w, str)
}
func (o *gql) WriteWebsocketMessage(mt int, socketId string, requestID RequestID, response *HttpResponse) {
	r := `{"id":"` + string(requestID) + `","type":"next","payload":` + response.Data + `}`
	gqlHttp.WriteWebsocketMessage(mt, socketId, []byte(r))
}
