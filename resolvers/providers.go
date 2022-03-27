package resolvers

import (
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type ProviderData struct{
	IsAuthenticated bool `gql:"name=isAuthenticated"`
}
type Provider struct{
}

func NewProvider() resolvers.ObjectTypeInterface{
	var o resolvers.ObjectTypeInterface
	o = &Provider{};
	return o;
}
func(o *Provider) Resolver(resolver string, args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList,typename string)( r resolvers.DataReturn ){
	x := ProviderData{IsAuthenticated:true};
	r = x;
	return r;
}
