package objectTypes

import (
	"github.com/pjmd89/gogql/lib/gql/introspection"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gogql/lib/gql/resolvers/directives"
	"github.com/vektah/gqlparser/v2/ast"
)

type Enum struct{
	schema resolvers.Schema
	directives map[string]resolvers.Directive
}

func NewEnumValue(schema resolvers.Schema,directives map[string]resolvers.Directive) resolvers.ObjectTypeInterface{
	var _type resolvers.ObjectTypeInterface
	_type = &Enum{schema:schema, directives: directives};

	return _type;
}

func(o *Enum) Resolver(resolver string, args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList,typename string) ( r resolvers.DataReturn ){
	
	switch(resolver){
		case "enumValues":
			r = o.enumValues(args, parent, directives, typename);
			break;
		default:
	}
	
	return r;
}
func(o *Enum) enumValues(args resolvers.Args, parent resolvers.Parent, directiveList resolvers.DirectiveList, typename string) (r resolvers.DataReturn ){	
	thisParent := parent.(introspection.Type);
	includeDeprecated := false;
	if(args["includeDeprecated"] != nil){
		includeDeprecated = args["includeDeprecated"].(bool);
	}
	switch(thisParent.Kind){
	case introspection.TYPEKIND_ENUM:
		r = make([]interface{},0);
		thisEnum := o.schema.Types[*thisParent.Name];
		for _,value := range thisEnum.EnumValues{
			x := introspection.EnumValue{};
			x.Name = value.Name;
			deprecatedResult := o.setDeprecate(value, thisParent);
			x.IsDeprecated =deprecatedResult.IsDeprecated;
			x.DeprecationReason =deprecatedResult.DeprecationReason;
			if value.Description != ""{
				x.Description = &value.Description;
			}
			if x.IsDeprecated == false{
				r = append(r.([]interface{}),x);
			}
			if includeDeprecated == true && x.IsDeprecated == true{
				r = append(r.([]interface{}),x);
			}
		}
	}
	return r;
}
func(o *Enum) setDeprecate(value *ast.EnumValueDefinition,thisParent introspection.Type) *directives.DeprecatedData{
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