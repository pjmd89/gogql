package main

import (
	"github.com/pjmd89/gogql/lib/gql"
	"github.com/pjmd89/gogql/lib/http"
	"github.com/pjmd89/gogql/resolvers"
)

func main(){

	schema := gql.Init("localhost","schema");
	schema.ObjectType("Provider",resolvers.NewProvider());
	myHttp := http.Init(schema);
	myHttp.Start();

}