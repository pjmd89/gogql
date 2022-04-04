package objectTypes

import (
	"github.com/pjmd89/gogql/lib/gql/introspection"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gogql/lib/gql/resolvers/directives"
	"github.com/pjmd89/gqlparser/v2/ast"
)
type Field struct{
	schema resolvers.Schema
	directives map[string]resolvers.Directive
}

func NewField(schema resolvers.Schema,directives map[string]resolvers.Directive) resolvers.ObjectTypeInterface{
	var _type resolvers.ObjectTypeInterface
	_type = &Field{schema:schema, directives: directives};

	return _type;
}
func(o *Field) Subscribe(info resolvers.ResolverInfo) ( r bool){
	return r;
}
func(o *Field) Resolver(info resolvers.ResolverInfo) ( r resolvers.DataReturn ){
	
	switch(info.Resolver){
		case "fields":
			r = o.fields(info.Args, info.Parent);
			break;
		default:
	}
	
	return r;
}

func(o *Field) fields(args resolvers.Args, parent resolvers.Parent) (r resolvers.DataReturn ){
	thisParent := parent.(introspection.Type);
	r = nil;
	includeDeprecated := false;
	if(args["includeDeprecated"] != nil){
		includeDeprecated = args["includeDeprecated"].(bool);
	}
	switch(thisParent.Kind){
	case introspection.TYPEKIND_OBJECT,introspection.TYPEKIND_INTERFACE:
		r = make([]interface{},0);
		for _,value:=range *thisParent.Fields{
			if value.Name != "__schema" && value.Name != "__type"{
				
				x := introspection.Field{};
				x.Args = value.Arguments;
				deprecatedResult := o.setDeprecate(value, thisParent);
				x.IsDeprecated =deprecatedResult.IsDeprecated;
				x.DeprecationReason =deprecatedResult.DeprecationReason;
				x.Name = value.Name;
				if value.Description != ""{
					x.Description = &value.Description;
				}
				x.Args = value.Arguments;
				x.Type = value.Type;
				if x.IsDeprecated == false{
					r = append(r.([]interface{}),x);
				}
				if includeDeprecated == true && x.IsDeprecated == true{
					r = append(r.([]interface{}),x);
				}
			}
		}
	}
	return r;
}
func(o *Field) setDeprecate(value *ast.FieldDefinition,thisParent introspection.Type) *directives.DeprecatedData{
	var deprecateDirectiveResult directives.DeprecatedData;
	if value.Directives != nil{
		for _,directive:=range value.Directives{
			switch directive.Name{
				case "deprecated":
					deprecateDirectiveResult = o.directives[directive.Name].Invoke(map[string]interface{}{},*thisParent.Name,value.Name).(directives.DeprecatedData);
			}
		}
	}
	return &deprecateDirectiveResult;
}
