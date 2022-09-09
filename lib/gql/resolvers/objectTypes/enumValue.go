package objectTypes

import (
	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/introspection"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gogql/lib/gql/resolvers/directives"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type Enum struct {
	schema     resolvers.Schema
	directives map[string]resolvers.Directive
}

func NewEnumValue(schema resolvers.Schema, directives map[string]resolvers.Directive) resolvers.ObjectTypeInterface {
	var _type resolvers.ObjectTypeInterface
	_type = &Enum{schema: schema, directives: directives}

	return _type
}
func (o *Enum) Subscribe(info resolvers.ResolverInfo) (r bool) {
	return r
}
func (o *Enum) Resolver(info resolvers.ResolverInfo) (r resolvers.DataReturn, err definitionError.GQLError) {

	switch info.Resolver {
	case "enumValues":
		r = o.enumValues(info.Args, info.Parent)
		break
	default:
	}

	return r, err
}
func (o *Enum) enumValues(args resolvers.Args, parent resolvers.Parent) (r resolvers.DataReturn) {
	thisParent := parent.(introspection.Type)
	includeDeprecated := false
	if args["includeDeprecated"] != nil {
		includeDeprecated = args["includeDeprecated"].(bool)
	}
	switch thisParent.Kind {
	case introspection.TYPEKIND_ENUM:
		r = make([]interface{}, 0)
		thisEnum := o.schema.Types[*thisParent.Name]
		for _, value := range thisEnum.EnumValues {
			x := introspection.EnumValue{}
			x.Name = value.Name
			deprecatedResult := o.setDeprecate(value, thisParent)
			x.IsDeprecated = deprecatedResult.IsDeprecated
			x.DeprecationReason = deprecatedResult.DeprecationReason
			if value.Description != "" {
				x.Description = &value.Description
			}
			if x.IsDeprecated == false {
				r = append(r.([]interface{}), x)
			}
			if includeDeprecated == true && x.IsDeprecated == true {
				r = append(r.([]interface{}), x)
			}
		}
	}
	return r
}
func (o *Enum) setDeprecate(value *ast.EnumValueDefinition, thisParent introspection.Type) *directives.DeprecatedData {
	var deprecateDirectiveResult directives.DeprecatedData
	if value.Directives != nil {
		for _, directive := range value.Directives {
			switch directive.Name {
			case "deprecated":
				deprecateDirectiveResults, _ := o.directives[directive.Name].Invoke(map[string]interface{}{}, *thisParent.Name, value.Name)
				deprecateDirectiveResult = deprecateDirectiveResults.(directives.DeprecatedData)
			}
		}
	}
	return &deprecateDirectiveResult
}
