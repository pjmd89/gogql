package gql

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)
func(o *gql) arguments(defaultArguments map[string]DefaultArguments,inputArguments map[string]interface{}) map[string]interface{}{
	args := make(map[string]interface{})

	return args
}
func (o *gql) parseArguments(argsRaw ast.ArgumentList, argsDefinition ast.ArgumentDefinitionList ) map[string]interface{} {
	args := make(map[string]interface{})
	inputArguments := o.parseInputArguments(argsRaw);
	defaultArguments := o.parseDefaultArguments(argsDefinition);
	o.arguments(defaultArguments,inputArguments);

	return args
}

func (o *gql) parseInputArguments(rawArgs ast.ArgumentList) map[string]interface{} {
	args := make(map[string]interface{},0);
	if len(rawArgs) > 0{
		for _,vArgs := range rawArgs{
			if len(vArgs.Value.Children) > 0{
				args[vArgs.Name] = o.parseArgChildren(vArgs.Value.Children);
			}else{
				switch(vArgs.Value.Definition.Kind){
				case "INPUT_OBJECT":
					args[vArgs.Name] = o.variables[vArgs.Value.Raw]
				default:
					args[vArgs.Name] = vArgs.Value.Raw;
				}
			}
		}
	}
	fmt.Println(args);
	return args;
}

func (o *gql) parseArgChildren(rawArgs ast.ChildValueList) interface{} {
	var args interface{}
	mapArgs := make(map[string]interface{},0);
	sliceArgs := make([]interface{},0);
	if len(rawArgs) > 0{
		for _,vArgs := range rawArgs{
			if vArgs.Name != ""{
				mapArgs[vArgs.Name] = vArgs.Value.Raw
				if len(vArgs.Value.Children) > 0 {
					mapArgs[vArgs.Name] = o.parseArgChildren(vArgs.Value.Children)
				}
			}else{
				if len(vArgs.Value.Children) > 0 {
					sliceArgs = append(sliceArgs,o.parseArgChildren(vArgs.Value.Children));
				} else {
					sliceArgs = append(sliceArgs,vArgs.Value.Raw);
				}
			}
		}
	}
	if len(mapArgs) > 0{
		args = mapArgs
	}else{
		args = sliceArgs
	}
	return args;
}
func(o *gql) parseDefaultArguments(argsDefinition ast.ArgumentDefinitionList ) map[string]DefaultArguments {
	args := make(map[string]DefaultArguments)
	if len(argsDefinition) > 0{
		for _,dArgs := range argsDefinition{
			defaultArgument := DefaultArguments{};
			defaultArgument.Name = dArgs.Name;
			defaultArgument.NonNull = dArgs.Type.NonNull;
			defaultArgument.Type = dArgs.Type.NamedType;
			if dArgs.DefaultValue != nil{
				defaultArgument.Value = dArgs.DefaultValue.Raw;
				if dArgs.DefaultValue.Children != nil{
					defaultArgument.Value = o.parseArgChildren(dArgs.DefaultValue.Children);
				}
			}
			switch dArgs.Type.NamedType {
			case "","ID","String","Int","Float","Boolean":
			default:
				argType := o.schema.Types[dArgs.Type.NamedType];
				switch argType.Kind {
				case "INPUT_OBJECT":
					defaultArgument.Value = o.parseFields(o.schema.Types[dArgs.Type.NamedType].Fields);
				case "ENUM":
				case "SCALAR":
				default:
				}
				fmt.Println(argType)
			}
			args[dArgs.Name] = defaultArgument;
		}
	}
	return args;
}
func(o *gql) parseFields(fields ast.FieldList) map[string]DefaultArguments{
	args := make(map[string]DefaultArguments,0);
	for _,field := range fields{
		argField := DefaultArguments{};
		argField.Name = field.Name;
		argField.NonNull = field.Type.NonNull;
		argField.Type = field.Type.NamedType;
		//argField.Kind = field.Type.;
		if field.DefaultValue != nil{
			argField.Value = field.DefaultValue.Raw;
			if field.DefaultValue.Children != nil{
				argField.Value = o.parseArgChildren(field.DefaultValue.Children);
			}
		}
		switch field.Type.NamedType {
		case "","ID","String","Int","Float","Boolean":
		default:
			argType := o.schema.Types[field.Type.NamedType];
			switch argType.Kind {
			case "INPUT_OBJECT":
				argField.Value = o.parseFields(o.schema.Types[field.Type.NamedType].Fields);
			case "ENUM":
			case "SCALAR":
			default:
			}
			fmt.Println(argType)
		}
		args[field.Name] = argField;
	}
	return args
}