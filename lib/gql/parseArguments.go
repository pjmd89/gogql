package gql

import (
	"log"

	"github.com/pjmd89/gqlparser/v2/ast"
)

/*
1. setear todos los valores por defecto
2. setear los valores que vienen del input
3. verificar los valores que sean array para parsearlos
4. validar el tipo de datos de cada argumento
*/
type varDef interface {
	*ast.Argument | *ast.ChildValue
}
type variableDef[T varDef] struct {
	data T
}

func (o *gql) parseArguments(argsRaw ast.ArgumentList, argsDefinition ast.ArgumentDefinitionList) map[string]interface{} {
	args := make(map[string]*DefaultArguments)
	for _, val := range argsDefinition {
		if args[val.Name] == nil {
			arg := &DefaultArguments{}
			arg.NonNull = val.Type.NonNull
			arg.Name = val.Name
			arg.Type = val.Type.NamedType
			if val.Type.Elem != nil {
				arg.Type = val.Type.Elem.NamedType
				arg.IsArray = true
			}
			arg.Kind = string(o.schema.Types[arg.Type].Kind)
			if val.DefaultValue != nil {
				arg.Value = val.DefaultValue.Raw
			}
			args[val.Name] = arg
		}
	}
	for _, argRaw := range argsRaw {
		switch argRaw.Value.Kind {
		case 0:
			if argRaw.Value.VariableDefinition != nil {
				args[argRaw.Name].Value = o.variables[argRaw.Value.Raw]
			}
			args[argRaw.Name].Value = o.setValue(argRaw)
		case 8, 9:
			args[argRaw.Name].Value = o.parseArgChildren(argRaw.Value.Children)
		default:
			args[argRaw.Name].Value = o.setValue(argRaw)
		}
	}
	r := o.validateArguments(args)
	return r
}
func (o *gql) validateArguments(args map[string]*DefaultArguments) map[string]interface{} {
	r := make(map[string]interface{}, 0)
	for k, v := range args {
		switch v.Kind {
		case "INPUT_OBJECT":
			r[k] = o.parseInputObject(v)
		case "SCALAR":
			if o.scalars[v.Type] != nil {
				r[k], _ = o.scalars[v.Type].Set(v.Value)
			} else {
				r[k] = v.Value
				log.Println("Scalar not found: ", v.Type)
			}
		default:
			r[k] = v.Value
		}
	}
	return r
}
func (o *gql) parseInputObject(argInput *DefaultArguments) (r interface{}) {
	if argInput != nil {
		args := make(map[string]*DefaultArguments)
		inputObject := o.schema.Types[argInput.Type]
		for _, val := range inputObject.Fields {
			arg := &DefaultArguments{}
			arg.NonNull = val.Type.NonNull
			arg.Name = val.Name
			arg.Type = val.Type.NamedType
			if val.Type.Elem != nil {
				arg.Type = val.Type.Elem.NamedType
				arg.IsArray = true
			}
			arg.Kind = string(o.schema.Types[arg.Type].Kind)
			if val.DefaultValue != nil {
				arg.Value = val.DefaultValue.Raw
			}
			args[val.Name] = arg
		}
		if argInput.IsArray {
			re := make([]interface{}, 0)
			newArgs := make(map[string]*DefaultArguments, 0)
			for k, v := range args {
				var x *DefaultArguments = &DefaultArguments{}
				*x = *v
				newArgs[k] = x
			}
			if argInput.Value != nil {
				for _, v := range argInput.Value.([]interface{}) {
					for name, val := range v.(map[string]interface{}) {
						if args[name] != nil {
							newArgs[name].Value = val
						}
					}
					re = append(re, o.validateArguments(newArgs))
				}
			}
			r = re
		} else {
			if argInput.Value != nil {
				for k, v := range argInput.Value.(map[string]interface{}) {
					args[k].Value = v
				}
				r = o.validateArguments(args)
			} else {
				re := make(map[string]interface{}, 0)
				for k, v := range args {
					re[k] = v.Value
				}
				r = re
			}
		}
	}
	return r
}
func (o *gql) validateScalar(arg DefaultArguments) {

}
func (o *gql) parseArgChildren(rawArgs ast.ChildValueList) interface{} {
	var args interface{}
	mapArgs := make(map[string]interface{}, 0)
	sliceArgs := make([]interface{}, 0)
	if len(rawArgs) > 0 {
		for _, vArgs := range rawArgs {
			if vArgs.Name != "" {
				mapArgs[vArgs.Name] = o.setValue(vArgs)
				if len(vArgs.Value.Children) > 0 {
					mapArgs[vArgs.Name] = o.parseArgChildren(vArgs.Value.Children)
				}
			} else {
				if len(vArgs.Value.Children) > 0 {
					sliceArgs = append(sliceArgs, o.parseArgChildren(vArgs.Value.Children))
				} else {
					sliceArgs = append(sliceArgs, o.setValue(vArgs))
				}
			}
		}
	}
	if len(mapArgs) > 0 {
		args = mapArgs
	} else {
		args = sliceArgs
	}
	return args
}
func (o *gql) setValue(vArgs any) (r any) {
	switch vArgs.(type) {
	case *ast.ChildValue:
		nArgs := vArgs.(*ast.ChildValue)
		r = nArgs.Value.Raw
		if nArgs.Value.VariableDefinition != nil {
			r = o.typedValue(nArgs.Value.Raw, nArgs.Value.VariableDefinition.Type.NamedType)
		}
	case *ast.Argument:
		nArgs := vArgs.(*ast.Argument)
		r = nArgs.Value.Raw
		if nArgs.Value.VariableDefinition != nil && nArgs.Value.VariableDefinition.Definition.Kind == "SCALAR" {
			r = o.typedValue(nArgs.Value.Raw, nArgs.Value.VariableDefinition.Type.NamedType)
		} else {
			if o.variables[nArgs.Value.Raw] != nil {
				r = o.variables[nArgs.Value.Raw]
			}
		}
	}
	return
}
func (o *gql) typedValue(name string, typed string) (r interface{}) {
	r, _ = o.scalars[typed].Set(o.variables[name])
	return
}
