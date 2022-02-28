package main

import (
	"github.com/pjmd89/gogql/lib/gql"
	"github.com/pjmd89/gogql/lib/http"
)

///*
func main(){

	schema := gql.Init("localhost","schema");
	myHttp := http.Init(schema);
	myHttp.Start();

}
//*/
/*
//basic example whitout gql
func main(){
	gqlDefault := new(http.GQLDefault);
	myHttp := http.Init(gqlDefault);
	myHttp.Start();

}
//*/