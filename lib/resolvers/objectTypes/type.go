package objectTypes

import (
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/introspection"
	"github.com/pjmd89/gogql/lib/resolvers"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type Type struct {
	schema     resolvers.Schema
	directives map[string]resolvers.Directive
}

func NewType(schema resolvers.Schema, directives map[string]resolvers.Directive) resolvers.ObjectTypeInterface {
	var _type resolvers.ObjectTypeInterface
	_type = &Type{schema: schema, directives: directives}
	return _type
}
func (o *Type) Subscribe(info resolvers.ResolverInfo) (r bool) {
	return r
}
func (o *Type) Resolver(info resolvers.ResolverInfo) (r resolvers.DataReturn, err definitionError.GQLError) {

	switch info.Resolver {
	case "__type":
		r = o.__type(info.Args, info.Parent)
	case "types":
		r = o.types(info.Parent)
	case "type":
		r = o._type(info.Parent)
	case "interfaces":
		r = o.interfaces(info.Parent)
	case "possibleTypes":
		r = o.possibleTypes(info.Parent)
	case "ofType":
		r = o.ofType(info.Parent)
	case "queryType":
		r = o.queryType(info.Parent)
	case "mutationType":
		r = o.mutationType(info.Parent)
	case "subscriptionType":
		r = o.subscriptionType(info.Parent)
	default:
	}

	return r, err
}
func (o *Type) subscriptionType(parent resolvers.Parent) (r resolvers.DataReturn) {
	parentInfo := parent.(introspection.Schema)
	thisType := parentInfo.SubscriptionType
	if thisType != nil {
		x := introspection.Type{}
		x.Kind = *introspection.SetTypeKind(string(thisType.Kind))
		if thisType.Name != "" {
			x.Name = &thisType.Name
		}
		if thisType.Description != "" {
			x.Description = &thisType.Description
		}
		r = x
	}
	return r
}
func (o *Type) mutationType(parent resolvers.Parent) (r resolvers.DataReturn) {
	parentInfo := parent.(introspection.Schema)
	thisType := parentInfo.MutationType
	if thisType != nil {
		x := introspection.Type{}
		x.Kind = *introspection.SetTypeKind(string(thisType.Kind))
		if thisType.Name != "" {
			x.Name = &thisType.Name
		}
		if thisType.Description != "" {
			x.Description = &thisType.Description
		}
		r = x
	}
	return r
}
func (o *Type) queryType(parent resolvers.Parent) (r resolvers.DataReturn) {
	parentInfo := parent.(introspection.Schema)
	thisType := parentInfo.QueryType
	if thisType != nil {
		x := introspection.Type{}
		x.Kind = *introspection.SetTypeKind(string(thisType.Kind))
		if thisType.Name != "" {
			x.Name = &thisType.Name
		}
		if thisType.Description != "" {
			x.Description = &thisType.Description
		}
		x.Fields = &thisType.Fields
		x.Interfaces = &thisType.Interfaces
		r = x
	}
	return r
}
func (o *Type) ofType(parent resolvers.Parent) (r resolvers.DataReturn) {
	thisParent := parent.(introspection.Type)
	thisType := thisParent.OfType
	switch thisParent.Kind {
	case introspection.TYPEKIND_LIST:
		x := introspection.Type{}
		if thisType.NonNull == true {
			x.Kind = *introspection.SetTypeKind("NON_NULL")
			ofType := &ast.Type{}
			copier.Copy(&ofType, thisType)
			ofType.NonNull = false
			x.OfType = ofType
		} else {
			x.Name = &thisType.NamedType
			x.Kind = *introspection.SetTypeKind(string(o.schema.Types[thisType.NamedType].Kind))
		}
		r = x
	case introspection.TYPEKIND_NON_NULL:
		x := introspection.Type{}
		if thisType.Elem != nil {
			x.Kind = *introspection.SetTypeKind("LIST")
			x.OfType = thisType.Elem
		} else {
			x.Name = &thisType.NamedType
			x.Kind = *introspection.SetTypeKind(string(o.schema.Types[thisType.NamedType].Kind))
		}
		r = x
	}
	return r
}
func (o *Type) _type(parent resolvers.Parent) (r resolvers.DataReturn) {
	rValue := reflect.ValueOf(parent)
	switch rValue.Type() {
	case reflect.TypeOf(introspection.Field{}):
		thisParent := parent.(introspection.Field)
		thisType := thisParent.Type
		x := introspection.Type{}
		if thisType.NonNull == true {
			x.Kind = *introspection.SetTypeKind("NON_NULL")
			//var ofType *ast.Type;
			ofType := &ast.Type{}
			copier.Copy(&ofType, thisType)
			ofType.NonNull = false
			x.OfType = ofType
		}
		if thisType.Elem != nil && thisType.NonNull == false {
			x.Kind = *introspection.SetTypeKind("LIST")
			x.OfType = thisType.Elem
		}
		if thisType.NonNull == false && thisType.Elem == nil {
			x.Name = &thisType.NamedType
			x.Kind = *introspection.SetTypeKind(string(o.schema.Types[thisType.NamedType].Kind))
		}
		r = x
	case reflect.TypeOf(introspection.InputValue{}):
		thisParent := parent.(introspection.InputValue)
		thisType := thisParent.Type
		x := introspection.Type{}
		if thisType.NonNull == true {
			x.Kind = *introspection.SetTypeKind("NON_NULL")
			//var ofType *ast.Type;
			ofType := &ast.Type{}
			copier.Copy(&ofType, thisType)
			ofType.NonNull = false
			x.OfType = ofType
		}
		if thisType.Elem != nil && thisType.NonNull == false {
			x.Kind = *introspection.SetTypeKind("LIST")
			x.OfType = thisType.Elem
		}
		if thisType.NonNull == false && thisType.Elem == nil {
			x.Name = &thisType.NamedType
			x.Kind = *introspection.SetTypeKind(string(o.schema.Types[thisType.NamedType].Kind))
		}
		r = x
	}
	return r
}
func (o *Type) possibleTypes(parent resolvers.Parent) (r resolvers.DataReturn) {
	thisParent := parent.(introspection.Type)
	switch thisParent.Kind {
	case introspection.TYPEKIND_INTERFACE, introspection.TYPEKIND_UNION:
		r = make([]interface{}, 0)
		possibleTypes := o.schema.PossibleTypes[*thisParent.Name]
		for _, value := range possibleTypes {
			x := introspection.Type{}
			x.Kind = *introspection.SetTypeKind(string(value.Kind))
			if value.Name != "" {
				x.Name = &value.Name
			}
			if value.Description != "" {
				x.Description = &value.Description
			}
			x.Fields = &value.Fields
			x.Interfaces = &value.Interfaces
			r = append(r.([]interface{}), x)
		}
	}
	return r
}

func (o *Type) interfaces(parent resolvers.Parent) (r resolvers.DataReturn) {
	thisParent := parent.(introspection.Type)
	switch thisParent.Kind {
	case introspection.TYPEKIND_OBJECT:
		r = make([]interface{}, 0)
		if thisParent.Interfaces != nil {
			for _, interfaceType := range *thisParent.Interfaces {
				x := introspection.Type{}
				findType := o.schema.Types[interfaceType]
				x.Kind = *introspection.SetTypeKind(string(findType.Kind))
				if findType.Name != "" {
					x.Name = &findType.Name
				}
				if findType.Description != "" {
					x.Description = &findType.Description
				}
				x.Fields = &findType.Fields
				x.Interfaces = &findType.Interfaces
				r = append(r.([]interface{}), x)
			}
		}
	}
	return r
}
func (o *Type) __type(args resolvers.Args, parent resolvers.Parent) (r resolvers.DataReturn) {
	findType := o.schema.Types[args["name"].(string)]
	if findType != nil {
		x := introspection.Type{}
		x.Kind = *introspection.SetTypeKind(string(findType.Kind))
		if findType.Name != "" {
			x.Name = &findType.Name
		}
		if findType.Description != "" {
			x.Description = &findType.Description
		}
		x.Fields = &findType.Fields
		x.Interfaces = &findType.Interfaces
		r = x
	}
	return r
}

func (o *Type) types(parent resolvers.Parent) (r resolvers.DataReturn) {
	r = make([]interface{}, 0)
	for _, findType := range parent.(introspection.Schema).Types {
		x := introspection.Type{}
		x.Kind = *introspection.SetTypeKind(string(findType.Kind))
		if findType.Name != "" {
			x.Name = &findType.Name
		}
		if findType.Description != "" {
			x.Description = &findType.Description
		}
		x.Fields = &findType.Fields
		x.Interfaces = &findType.Interfaces
		r = append(r.([]interface{}), x)
	}
	return r
}
