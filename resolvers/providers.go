package resolvers

import (
	"github.com/pjmd89/gogql/lib/gql"
	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gogql/lib/http"
	"github.com/pjmd89/gogql/subscriptions"
)

type ProviderData struct{
	Name string `gql:"name=name"`
}
type Provider struct{}

func NewProvider() resolvers.ObjectTypeInterface{
	var o resolvers.ObjectTypeInterface
	o = &Provider{};
	return o;
}
func(o *Provider) Subscribe(info resolvers.ResolverInfo) ( r bool  ){
	switch info.Resolver{
	case "providerAdded":
		if http.Session.Get()["id"] != nil && info.SubscriptionValue!= nil{
			s := info.SubscriptionValue.(resolvers.DataReturn).(ProviderData);
			if s.Name == "s"{
				r = true;
			}
		}
		
	}
	return r;
}
func(o *Provider) Resolver(info resolvers.ResolverInfo )( r resolvers.DataReturn, err definitionError.Error ){
	switch info.Operation{
	case "query", "mutation":
		switch info.Resolver{
		case "createProvider":
			r = o.createProvider(info.Args,info.Directives,info.TypeName);
		}
	case "subscription":
		switch info.Resolver{
		case "providerAdded":
			r = o.providerAdded(info.Args,info.Directives,info.TypeName);
		}
	}
	return r, err;
}
func(o *Provider) createProvider(args resolvers.Args, directives resolvers.DirectiveList, typename string) ( r resolvers.DataReturn ){
	input := args["input"].(map[string]interface{});
	r = ProviderData{
		Name:input["name"].(string),
	};
	cookie := make(map[interface{}]interface{});
	cookie["name"] = r.(ProviderData).Name;
	cookie["id"] = 1;
	http.Session.Set(cookie);
	gql.PubSub.Publish(subscriptions.PROVIDER_ADDED,r)
	return r;
}
func(o *Provider) providerAdded(args resolvers.Args, directives resolvers.DirectiveList, typename string) ( r resolvers.DataReturn ){
	rx := &ProviderData{};
	cookie := http.Session.Get();
	rx.Name = cookie["name"].(string);
	r = rx;
	return r;
}
