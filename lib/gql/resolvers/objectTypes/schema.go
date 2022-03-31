package objectTypes

import (
	"github.com/pjmd89/gogql/lib/gql/introspection"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type schema struct{
	schema resolvers.Schema
	directives map[string]resolvers.Directive
}

func NewSchema(gqlSchema resolvers.Schema, directives map[string]resolvers.Directive)resolvers.ObjectTypeInterface{
	var schema_ resolvers.ObjectTypeInterface
	schema_ = &schema{schema:gqlSchema, directives: directives};
	return schema_;
}
func(o *schema) Subscribe(info resolvers.ResolverInfo) ( r bool, s resolvers.Subscription ){
	return r, s;
}
func(o *schema) Resolver(info resolvers.ResolverInfo) ( r resolvers.DataReturn ){

	switch(info.Resolver){
		case "__schema":
			r = o.__schema();
			break;
		default:
	}
	
	return r;
}
func(o *schema) __schema() ( r resolvers.DataReturn ){
	x := introspection.Schema{};
	x.Types = o.schema.Types;
	x.QueryType = o.schema.Query;
	x.MutationType = o.schema.Mutation;
	x.SubscriptionType = o.schema.Subscription;
	x.Directives = o.schema.Directives;
	r = x;
	return r;
}


