package gql

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pjmd89/gogql/lib"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gogql/lib/gql/resolvers/directives"
	"github.com/pjmd89/gogql/lib/gql/resolvers/objectTypes"
	"github.com/pjmd89/gogql/lib/gql/resolvers/scalars"
	gqlHttp "github.com/pjmd89/gogql/lib/http"

	"github.com/pjmd89/gqlparser/v2"
	"github.com/pjmd89/gqlparser/v2/ast"
)

func Init(serverName string, path string) *gql{
	gql := &gql{};
	gql.serverName = serverName;
	gql.loadSchema(path);
    gql.objectTypes = make(map[string]resolvers.ObjectTypeInterface);
    gql.directives = make(map[string]resolvers.Directive);
    gql.scalars = make(map[string]resolvers.Scalar);

    gql.objectTypes["__Schema"] = objectTypes.NewSchema(gql.schema,gql.directives);
    gql.objectTypes["__Type"]   = objectTypes.NewType(gql.schema,gql.directives);
    gql.objectTypes["__Field"]  = objectTypes.NewField(gql.schema,gql.directives);
    gql.objectTypes["__EnumValue"]  = objectTypes.NewEnumValue(gql.schema,gql.directives);
    gql.objectTypes["__InputValue"]  = objectTypes.NewInputValue(gql.schema,gql.directives);
    gql.objectTypes["__Directive"]  = objectTypes.NewDirective(gql.schema);
    /*
    gql.directives["include"] = directives.NewInclude(gql.schema);
    gql.directives["skip"] = directives.NewSkip(gql.schema);
    //*/
    gql.directives["deprecated"] = directives.NewDeprecated(gql.schema);
    gql.scalars["ID"] = scalars.NewIDScalar();
    gql.scalars["Boolean"] = scalars.NewBoolScalar();
    gql.scalars["String"] = scalars.NewStringScalar();
    gql.scalars["Int"] = scalars.NewIntScalar();
    gql.scalars["Float"] = scalars.NewFloatScalar();
	return gql;
}
func(o *gql) ObjectType(resolver string, object resolvers.ObjectTypeInterface){
    o.objectTypes[resolver] = object;
}
func(o *gql) Directive(resolver string, object resolvers.Directive){
    o.directives[resolver] = object;
}
func(o *gql) Scalar(resolver string, object resolvers.Scalar){
    o.scalars[resolver] = object;
}
func(o *gql) loadSchema(path string){
	var schema []*ast.Source;
    files := lib.ScanDir(path);

    for _,file := range files{
        content, err := ioutil.ReadFile(file);
        if err != nil {
            log.Fatal(err);
        }
        text := string(content);
        schema = append(schema, &ast.Source{Name: file,Input: text,BuiltIn: true});
    }
    parser,err := gqlparser.LoadSchema(schema...);
    
    for _,subs := range parser.Subscription.Fields{
        PubSub.createSubscriptionEvent(OperationID(subs.Name));
    }
    if err !=nil{
        log.Fatal(err);
    }
    o.schema = parser;
}
func(o *gql) GQLRenderSubscription (mt int, message []byte , socketId string) {
    var request WebSocketRequest;
    json.Unmarshal(message,&request);
    //var response *HttpResponse = &HttpResponse{};
    switch request.Type {
    case "connection_init":
        r := `{"type":"connection_ack","payload":{}}`;
        gqlHttp.WriteWebsocketMessage(mt,socketId,[]byte(r));
    case "subscribe":
        o.WebsocketResponse(request.Payload,socketId,RequestID(request.Id), mt);
        //r = `{"id":"`+request.Id+`","type":"next","payload":`+response.Data+`}`;
    case "ping":
        r := `{"type":"pong","payload":{}}`;
        gqlHttp.WriteWebsocketMessage(mt,socketId,[]byte(r));
    case "complete":
        r := `{"id":"`+request.Id+`","type":"complete"}`;
        gqlHttp.WriteWebsocketMessage(mt,socketId,[]byte(r));
    default:
        
    }
}
func(o *gql) GQLRender(w http.ResponseWriter,r *http.Request) string{
    var request HttpRequest;
    json.NewDecoder(r.Body).Decode(&request)
    response := o.response(request);
    rx := response.Data;

    return rx;
}
func(o *gql) WriteWebsocketMessage(mt int,socketId string, requestID RequestID,response *HttpResponse){
    r := `{"id":"`+string(requestID)+`","type":"next","payload":`+response.Data+`}`;
    gqlHttp.WriteWebsocketMessage(mt,socketId,[]byte(r));
}
func(o *gql) GetServerName() string{
	return o.serverName;
}
