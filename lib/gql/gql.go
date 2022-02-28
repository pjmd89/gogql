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

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func Init(serverName string, path string) *gql{
	gql := &gql{};
	gql.serverName = serverName;
	gql.loadSchema(path);
    gql.objectTypes = make(map[string]resolvers.ObjectTypeInterface);
    gql.Directives = make(map[string]resolvers.Directive);

    gql.objectTypes["__Schema"] = objectTypes.NewSchema(gql.schema,gql.Directives);
    gql.objectTypes["__Type"]   = objectTypes.NewType(gql.schema,gql.Directives);
    gql.objectTypes["__Field"]  = objectTypes.NewField(gql.schema,gql.Directives);
    gql.objectTypes["__EnumValue"]  = objectTypes.NewEnumValue(gql.schema,gql.Directives);
    gql.objectTypes["__InputValue"]  = objectTypes.NewInputValue(gql.schema,gql.Directives);
    gql.objectTypes["__Directive"]  = objectTypes.NewDirective(gql.schema);
    /*
    gql.Directives["include"] = directives.NewInclude(gql.schema);
    gql.Directives["skip"] = directives.NewSkip(gql.schema);
    //*/
    gql.Directives["deprecated"] = directives.NewDeprecated(gql.schema);

    //anadir directivas y resolvers;
	return gql;
}
func(o *gql) ObjectType(resolver string, object resolvers.ObjectTypeInterface){
    o.objectTypes[resolver] = object;
}
func(o *gql) Directive(resolver string, object resolvers.Directive){
    o.Directives[resolver] = object;
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
    
    if err !=nil{
        log.Fatal(err);
    }
    o.schema = parser;
}
func(o *gql) GQLRender(w http.ResponseWriter,r *http.Request) string{
    var request HttpRequest;
    json.NewDecoder(r.Body).Decode(&request)
    response := o.response(request);
    rx := response.Data;

    return rx;
}
func(o *gql) GetServerName() string{
	return o.serverName;
}
