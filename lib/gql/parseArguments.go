package gql

import (
	"log"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
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

func (o *gql) parseArguments(argsRaw ast.ArgumentList, argsDefinition ast.ArgumentDefinitionList, vars map[string]any) (r map[string]interface{}, err definitionError.GQLError) {
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
				args[argRaw.Name].Value = vars[argRaw.Value.Raw]
			}
			args[argRaw.Name].Value = o.setValue(argRaw, vars)
		case 8, 9:
			args[argRaw.Name].Value = o.parseArgChildren(argRaw.Value.Children, vars)
		default:
			args[argRaw.Name].Value = o.setValue(argRaw, vars)
		}
	}
	r, err = o.validateArguments(args)
	return
}
func (o *gql) validateArguments(args map[string]*DefaultArguments) (r map[string]interface{}, err definitionError.GQLError) {
	r = make(map[string]interface{}, 0)
	for k, v := range args {
		switch v.Kind {
		case "INPUT_OBJECT":
			r[k], err = o.parseInputObject(v)
			if err != nil {
				return
			}
		case "SCALAR":
			if o.scalars[v.Type] != nil {
				if v.IsArray && v.Value != nil {
					parsedValue := []any{}
					for _, vV := range v.Value.([]any) {
						vScalar, vError := o.scalars[v.Type].Set(vV)
						if vError != nil {
							r = make(map[string]interface{}, 0)
							err = vError
							return
						}
						parsedValue = append(parsedValue, vScalar)
					}
					r[k] = parsedValue
				} else {
					vScalar, vError := o.scalars[v.Type].Set(v.Value)
					if vError != nil {
						r = make(map[string]interface{}, 0)
						err = vError
						return
					}
					r[k] = vScalar
				}

			} else {
				r[k] = v.Value
				log.Println("Scalar not found: ", v.Type)
			}
		default:
			r[k] = v.Value
		}
	}
	return
}
func (o *gql) parseInputObject(argInput *DefaultArguments) (r interface{}, err definitionError.GQLError) {
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
			valueType := o.schema.Types[val.Type.NamedType]
			if val.DefaultValue != nil {
				arg.Value = val.DefaultValue.Raw
			}
			if valueType != nil && valueType.Kind == "SCALAR" {
				if arg.IsArray {
					parsedValue := []any{}
					for _, vV := range arg.Value.([]any) {
						vScalar, vError := o.scalars[arg.Type].Set(vV)
						if vError != nil {
							r = make(map[string]interface{}, 0)
							err = vError
							return
						}
						parsedValue = append(parsedValue, vScalar)
					}
					arg.Value = parsedValue
				} else {
					vScalar, vError := o.scalars[valueType.Name].Set(arg.Value)
					if vError != nil {
						r = make(map[string]interface{}, 0)
						err = vError
						return
					}
					arg.Value = vScalar
				}
			}
			args[val.Name] = arg
		}
		if argInput.IsArray {
			re := make([]interface{}, 0)
			if argInput.Value != nil {
				switch argInput.Value.(type) {
				case []any:
					for _, v := range argInput.Value.([]interface{}) {
						newArgs := make(map[string]*DefaultArguments, 0)
						for k, v := range args {
							var x *DefaultArguments = &DefaultArguments{}
							*x = *v
							newArgs[k] = x
						}
						for name, val := range v.(map[string]interface{}) {
							if args[name] != nil {
								newArgs[name].Value = val
							}
						}
						vValue, vError := o.validateArguments(newArgs)
						if vError != nil {
							return vValue, vError
						}
						re = append(re, vValue)
					}
				default:
					log.Printf("variable %s is not an array ", argInput.Name)
					re = nil
				}
			}
			r = re
		} else {
			if argInput.Value != nil {
				for k, v := range argInput.Value.(map[string]interface{}) {
					args[k].Value = v
				}
				vValue, vError := o.validateArguments(args)
				if vError != nil {
					return vValue, vError
				}
				r = vValue
			} else {
				re := make(map[string]interface{}, 0)
				for k, v := range args {
					re[k] = v.Value
				}
				r = re
			}
		}
	}
	return
}
func (o *gql) parseArgChildren(rawArgs ast.ChildValueList, vars map[string]any) interface{} {
	var args interface{}
	mapArgs := make(map[string]interface{}, 0)
	sliceArgs := make([]interface{}, 0)
	if len(rawArgs) > 0 {
		for _, vArgs := range rawArgs {
			if vArgs.Name != "" {
				mapArgs[vArgs.Name] = o.setValue(vArgs, vars)
				if len(vArgs.Value.Children) > 0 {
					mapArgs[vArgs.Name] = o.parseArgChildren(vArgs.Value.Children, vars)
				}
			} else {
				if len(vArgs.Value.Children) > 0 {
					sliceArgs = append(sliceArgs, o.parseArgChildren(vArgs.Value.Children, vars))
				} else {
					sliceArgs = append(sliceArgs, o.setValue(vArgs, vars))
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
func (o *gql) setValue(vArgs any, vars map[string]any) (r any) {
	switch vArgs.(type) {
	case *ast.ChildValue:
		nArgs := vArgs.(*ast.ChildValue)
		r = nArgs.Value.Raw
		if nArgs.Value.VariableDefinition != nil {
			r = vars[nArgs.Value.Raw]
			//r = o.typedValue(nArgs.Value.Raw, nArgs.Value.VariableDefinition.Type.NamedType)
		}
	case *ast.Argument:
		nArgs := vArgs.(*ast.Argument)
		r = nArgs.Value.Raw
		if nArgs.Value.VariableDefinition != nil && nArgs.Value.VariableDefinition.Definition.Kind == "SCALAR" {
			r = vars[nArgs.Value.Raw]
			//r = o.typedValue(nArgs.Value.Raw, nArgs.Value.VariableDefinition.Type.NamedType)
		} else {
			if vars[nArgs.Value.Raw] != nil {
				r = vars[nArgs.Value.Raw]
			}
		}
	}
	return
}
