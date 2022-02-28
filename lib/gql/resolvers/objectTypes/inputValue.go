package objectTypes

import (
	"reflect"

	"github.com/pjmd89/gogql/lib/gql/introspection"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/vektah/gqlparser/v2/ast"
)
type InputValue struct{
	schema resolvers.Schema
	directives map[string]resolvers.Directive
}

func NewInputValue(schema resolvers.Schema,directives map[string]resolvers.Directive) resolvers.ObjectTypeInterface{
	var _type resolvers.ObjectTypeInterface
	_type = &InputValue{schema:schema, directives: directives};

	return _type;
}
func(o *InputValue) Resolver(resolver string, args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList,typename string) ( r resolvers.DataReturn ){
	
	switch(resolver){
	case "args":
		r = o.args(args, parent, directives, typename);
	case "inputFields":
		r = o.fields(args, parent, directives, typename);
	default:
	}
	
	return r;
}
func(o *InputValue) args(args resolvers.Args, parent resolvers.Parent, directiveList resolvers.DirectiveList, typename string) (r resolvers.DataReturn ){
	r = make([]interface{},0);
	var arguments ast.ArgumentDefinitionList;
	rValue := reflect.ValueOf(parent);
	switch(rValue.Type()){
	case reflect.TypeOf(introspection.Field{}):
		arguments = (parent.(introspection.Field)).Args;
	case reflect.TypeOf(introspection.Directive{}):
		arguments = (parent.(introspection.Directive)).Args;
	default:
		return r;
	}
	for _,value:=range arguments{
		x := introspection.InputValue{};
		x.Name = value.Name;
		if value.Description != ""{
			x.Description = &value.Description;
		}
		if value.DefaultValue != nil{
			x.DefaultValue = &value.DefaultValue.Raw;	
		}
		x.Type = value.Type;
		r = append(r.([]interface{}),x);
	}
	return r;
}
func(o *InputValue) fields(args resolvers.Args, parent resolvers.Parent, directiveList resolvers.DirectiveList, typename string) (r resolvers.DataReturn ){
	thisParent := parent.(introspection.Type);
	switch(thisParent.Kind){
	case introspection.TYPEKIND_INPUT_OBJECT:
		if thisParent.Fields != nil{
			r = make([]interface{},0);
			for _,value:=range *thisParent.Fields{
				x := introspection.InputValue{};
				x.Name = value.Name;
				x.Description = &value.Description;
				if value.DefaultValue != nil{
					x.DefaultValue = &value.DefaultValue.Raw;	
				}
				x.Type = value.Type;
				r = append(r.([]interface{}),x);
			}
		}
	}
	return r;
}
