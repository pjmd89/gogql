package objectTypes

import (
	"github.com/pjmd89/gogql/lib/gql/introspection"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type Directive struct{
	definition resolvers.Definition
	schema resolvers.Schema
}

func NewDirective(schema resolvers.Schema) resolvers.ObjectTypeInterface{
	var _type resolvers.ObjectTypeInterface
	_type = &Directive{schema:schema};

	return _type;
}
func (o *Directive) SetDefinition(definition resolvers.Definition){
	o.definition = definition;
}
func(o *Directive) Resolver(resolver string, args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList,typename string) ( r resolvers.DataReturn ){
	thisParent := parent.(introspection.Schema);
	r = make([]interface{},0);
	for _,value := range thisParent.Directives{
		x:=introspection.Directive{};
		x.Name = value.Name;
		if value.Description != ""{
			x.Description = &value.Description;
		}
		x.Locations = value.Locations;
		x.Args = value.Arguments;
		r = append(r.([]interface{}),x);
	}
	return r;
}