package directives

import (
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gqlparser/v2/ast"
)

type Deprecated struct{
	schema resolvers.Schema
}
type DeprecatedData struct{
	IsDeprecated bool
	DeprecationReason *string
}
func NewDeprecated(schema resolvers.Schema) resolvers.Directive{
	var _type resolvers.Directive
	_type = &Deprecated{schema:schema};
	return _type;
	
}
func (o *Deprecated) Invoke(args map[string]interface{},typeName string, fieldName string, directiveDefinition *ast.DirectiveDefinition ) resolvers.DataReturn{
	r := DeprecatedData{};
	if o.schema.Types[typeName] != nil{
		switch(o.schema.Types[typeName].Kind){
		case "ENUM":
			o.parseEnumValues(o.schema.Types[typeName].EnumValues, fieldName, &r);
		default:
			o.parseFields(o.schema.Types[typeName].Fields, fieldName, &r);
		}
	}
	return r;

}

func(o *Deprecated) parseFields(fields ast.FieldList, fieldName string, r *DeprecatedData){
	for _,field := range fields{
		if field.Name == fieldName{
			if field.Directives != nil{
				for _,directive := range field.Directives{
					if directive.Name == "deprecated"{
						if directive.Arguments != nil{
							for _,arg := range directive.Arguments{
								if arg.Name == "reason"{
									r.IsDeprecated = true;
									str := arg.Value.Raw;
									r.DeprecationReason = &str;
									break;
								}
							}
						}
					}
				}
			}
			break;
		}
	}
}
func(o *Deprecated) parseEnumValues(fields ast.EnumValueList, fieldName string, r *DeprecatedData){
	for _,field := range fields{
		if field.Name == fieldName{
			if field.Directives != nil{
				for _,directive := range field.Directives{
					if directive.Name == "deprecated"{
						if directive.Arguments != nil{
							for _,arg := range directive.Arguments{
								if arg.Name == "reason"{
									r.IsDeprecated = true;
									str := arg.Value.Raw;
									r.DeprecationReason = &str;
									break;
								}
							}
						}
					}
				}
			}
			break;
		}
	}
}