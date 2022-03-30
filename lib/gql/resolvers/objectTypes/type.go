package objectTypes

import (
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/pjmd89/gogql/lib/gql/introspection"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type Type struct{
	schema resolvers.Schema
	directives map[string]resolvers.Directive
}

func NewType(schema resolvers.Schema,directives map[string]resolvers.Directive) resolvers.ObjectTypeInterface{
	var _type resolvers.ObjectTypeInterface
	_type = &Type{schema:schema, directives: directives};
	return _type;
}
func(o *Type) Resolver(resolver string, args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList, typename string) ( r resolvers.DataReturn ){
	
	switch(resolver){
		case "__type":
			r = o.__type(args, parent, directives, typename);
		case "types":
			r = o.types(args, parent, directives, typename);
		case "type":
			r = o._type(args, parent, directives, typename);
		case "interfaces":
			r = o.interfaces(args, parent, directives, typename);
		case "possibleTypes":
			r = o.possibleTypes(args, parent, directives, typename);
		case "ofType":
			r = o.ofType(args, parent, directives, typename);
		case "queryType":
			r = o.queryType(args, parent, directives, typename);
		case "mutationType":
			r = o.mutationType(args, parent, directives, typename);
		case "subscriptionType":
			r = o.subscriptionType(args, parent, directives, typename);
		default:
	}
	
	return r;
}
func(o *Type) subscriptionType(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList, typename string) (r resolvers.DataReturn ){
	parentInfo := parent.(introspection.Schema);
	thisType := parentInfo.SubscriptionType;
	if thisType != nil{
		x := introspection.Type{};
		x.Kind = *introspection.SetTypeKind(string(thisType.Kind));
		if thisType.Name != ""{
			x.Name = &thisType.Name;
		}
		if thisType.Description != ""{
			x.Description = &thisType.Description;
		}
		r = x;
	}
	return r;
}
func(o *Type) mutationType(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList, typename string) (r resolvers.DataReturn ){
	parentInfo := parent.(introspection.Schema);
	thisType := parentInfo.MutationType;
	if thisType != nil{
		x := introspection.Type{};
		x.Kind = *introspection.SetTypeKind(string(thisType.Kind));
		if thisType.Name != ""{
			x.Name = &thisType.Name;
		}
		if thisType.Description != ""{
			x.Description = &thisType.Description;
		}
		r = x;
	}
	return r;
}
func(o *Type) queryType(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList, typename string) (r resolvers.DataReturn ){
	parentInfo := parent.(introspection.Schema);
	thisType := parentInfo.QueryType;
	if thisType != nil{
		x := introspection.Type{};
		x.Kind = *introspection.SetTypeKind(string(thisType.Kind));
		if thisType.Name != ""{
			x.Name = &thisType.Name;
		}
		if thisType.Description != ""{
			x.Description = &thisType.Description;
		}
		x.Fields = &thisType.Fields;
		x.Interfaces =  &thisType.Interfaces;
		r = x;
	}
	return r;
}
func(o *Type) ofType(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList, typename string) ( r resolvers.DataReturn ){
	thisParent := parent.(introspection.Type);
	thisType := thisParent.OfType;
	switch(thisParent.Kind){
	case introspection.TYPEKIND_LIST:
		x := introspection.Type{};
		if thisType.NonNull == true{
			x.Kind = *introspection.SetTypeKind("NON_NULL");
			ofType := &ast.Type{};
			copier.Copy(&ofType,thisType);
			ofType.NonNull = false;
			x.OfType = ofType;
		} else {
			x.Name = &thisType.NamedType;
			x.Kind = *introspection.SetTypeKind(string(o.schema.Types[thisType.NamedType].Kind));
		}
		r = x;
	case introspection.TYPEKIND_NON_NULL:
		x := introspection.Type{};
		if thisType.Elem !=nil{
			x.Kind = *introspection.SetTypeKind("LIST");
			x.OfType = thisType.Elem;
		}else{
			x.Name = &thisType.NamedType;
			x.Kind = *introspection.SetTypeKind(string(o.schema.Types[thisType.NamedType].Kind));
		}
		r = x;
	}
	return r;
}
func(o *Type) _type(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList,typename string) (r resolvers.DataReturn ){
	rValue := reflect.ValueOf(parent);
	switch(rValue.Type()){
	case reflect.TypeOf(introspection.Field{}):
		thisParent := parent.(introspection.Field);
		thisType := thisParent.Type;
		x := introspection.Type{};
		if thisType.NonNull == true{
			x.Kind = *introspection.SetTypeKind("NON_NULL");
			//var ofType *ast.Type;
			ofType := &ast.Type{};
			copier.Copy(&ofType,thisType);
			ofType.NonNull = false;
			x.OfType = ofType;
		}
		if thisType.Elem !=nil && thisType.NonNull == false{
			x.Kind = *introspection.SetTypeKind("LIST");
			x.OfType = thisType.Elem;
		}
		if thisType.NonNull == false && thisType.Elem ==nil{
			x.Name = &thisType.NamedType;
			x.Kind = *introspection.SetTypeKind(string(o.schema.Types[thisType.NamedType].Kind));
		}
		r = x;
	case reflect.TypeOf(introspection.InputValue{}):
		thisParent := parent.(introspection.InputValue);
		thisType := thisParent.Type;
		x := introspection.Type{};
		if thisType.NonNull == true{
			x.Kind = *introspection.SetTypeKind("NON_NULL");
			//var ofType *ast.Type;
			ofType := &ast.Type{};
			copier.Copy(&ofType,thisType);
			ofType.NonNull = false;
			x.OfType = ofType;
		}
		if thisType.Elem !=nil && thisType.NonNull == false{
			x.Kind = *introspection.SetTypeKind("LIST");
			x.OfType = thisType.Elem;
		}
		if thisType.NonNull == false && thisType.Elem ==nil{
			x.Name = &thisType.NamedType;
			x.Kind = *introspection.SetTypeKind(string(o.schema.Types[thisType.NamedType].Kind));
		}
		r = x;
	}
	return r;
}
func(o *Type) possibleTypes(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList, typename string) ( r resolvers.DataReturn ){
	thisParent := parent.(introspection.Type);
	switch(thisParent.Kind){
		case introspection.TYPEKIND_INTERFACE, introspection.TYPEKIND_UNION:
		r = make([]interface{},0);
		possibleTypes := o.schema.PossibleTypes[*thisParent.Name];
		for _,value := range possibleTypes{
			x := introspection.Type{};
			x.Kind = *introspection.SetTypeKind(string(value.Kind));
			if value.Name != ""{
				x.Name = &value.Name;
			}
			if value.Description != ""{
				x.Description = &value.Description;
			}
			x.Fields = &value.Fields;
			x.Interfaces =  &value.Interfaces;
			r = append(r.([]interface{}),x);
		}
	}
	return r;
}

func(o *Type) interfaces(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList, typename string) ( r resolvers.DataReturn ){
	thisParent := parent.(introspection.Type);
	switch(thisParent.Kind){
		case introspection.TYPEKIND_OBJECT:
		r = make([]interface{},0);
		if thisParent.Interfaces != nil{
			for _, interfaceType := range *thisParent.Interfaces{
				x := introspection.Type{};
				findType:= o.schema.Types[interfaceType];
				x.Kind = *introspection.SetTypeKind(string(findType.Kind));
				if findType.Name != ""{
					x.Name = &findType.Name;
				}
				if findType.Description != ""{
					x.Description = &findType.Description;
				}
				x.Fields = &findType.Fields;
				x.Interfaces =  &findType.Interfaces;
				r = append(r.([]interface{}),x);
			}
		}
	}
	return r;
}
func(o *Type) __type(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList, typename string) (r resolvers.DataReturn ){
	findType:= o.schema.Types[args["name"].(string)];
	if findType != nil{
		x := introspection.Type{};
		x.Kind = *introspection.SetTypeKind(string(findType.Kind));
		if findType.Name != ""{
			x.Name = &findType.Name;
		}
		if findType.Description != ""{
			x.Description = &findType.Description;
		}
		x.Fields = &findType.Fields;
		x.Interfaces =  &findType.Interfaces;
		r = x;
	}
	return r;
}

func(o *Type) types(args resolvers.Args, parent resolvers.Parent, directives resolvers.DirectiveList,typename string) ( r resolvers.DataReturn ){
	r = make([]interface{},0);
	for _,findType := range parent.(introspection.Schema).Types{
		x := introspection.Type{};
		x.Kind = *introspection.SetTypeKind(string(findType.Kind));
		if findType.Name != ""{
			x.Name = &findType.Name;
		}
		if findType.Description != ""{
			x.Description = &findType.Description;
		}
		x.Fields = &findType.Fields;
		x.Interfaces =  &findType.Interfaces;
		r = append(r.([]interface{}),x);
	}
	return r;
}