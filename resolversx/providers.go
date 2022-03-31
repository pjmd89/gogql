package resolvers

import (
	"github.com/pjmd89/gogql/lib/gql/pubsub"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gogql/lib/http"
	"github.com/pjmd89/gogql/subscriptions"
)

type ProviderData struct{
	Name string `gql:"name=name"`
}
type Provider struct{
}

func NewProvider() resolvers.ObjectTypeInterface{
	var o resolvers.ObjectTypeInterface
	o = &Provider{};
	return o;
}
func(o *Provider) Subscribe(info resolvers.ResolverInfo) ( r bool, s resolvers.Subscription ){
	switch info.Resolver{
	case "providerAdded":
		s = pubsub.Storage.Subscribe(subscriptions.LOGGED);
		if http.Session.Get()["id"] != nil && s.Value != nil{
			r = true;
		}
	}
	return r, s;
}
func(o *Provider) Resolver(info resolvers.ResolverInfo )( r resolvers.DataReturn ){
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
	return r;
}
func(o *Provider) createProvider(args resolvers.Args, directives resolvers.DirectiveList, typename string) ( r resolvers.DataReturn ){
	input := args["input"].(map[string]interface{});
	r = ProviderData{
		Name:input["name"].(string),
	};
	cookie := make(map[interface{}]interface{});
	cookie["name"] = r.(ProviderData).Name;
	cookie["id"] = 1;
	http.Session.New(cookie);
	pubsub.Storage.Publish(subscriptions.LOGGED,r);
	return r;
}
func(o *Provider) providerAdded(args resolvers.Args, directives resolvers.DirectiveList, typename string) ( r resolvers.DataReturn ){
	rx := &ProviderData{};
	cookie := http.Session.Get();
	rx.Name = cookie["name"].(string);
	r = rx;
	return r;
}
