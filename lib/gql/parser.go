package gql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator"
)

type Response interface{}
func(o *gql) response(request HttpRequest) *HttpResponse{
	response := &HttpResponse{};
	document,err := gqlparser.LoadQuery(o.schema,request.Query);
	
	if err != nil{
		fmt.Println(err.Error());

	}
	if document != nil{
		parse := document.Operations;
		if strings.Trim(request.OperationName," ") != ""{
			for _,operation := range document.Operations{
				if operation.Name == request.OperationName{
					parse = ast.OperationList{operation};
					break;
				}
			}
		}
		rx := o.operationParse(parse,request.Variables)
		response.Data = fmt.Sprintf("%v",rx["data"]);
	}
	return response;
}

func(o *gql) operationParse(parse ast.OperationList, variables map[string]interface{}) (map[string]interface{}){
	//este es el nombre de la query si es que lo tiene, si no, ejecuta la query unica
	prepareToSend := make(map[string]interface{},0);
	for _,operation :=range parse{
		o.setVariables(o.schema,operation,variables);
		var dataReturn resolvers.DataReturn;
		data := make(map[string]interface{},0);
		data["data"] = o.selectionSetParse(operation.SelectionSet, dataReturn ,dataReturn,nil);
		prepareToSend["data"] = o.jsonResponse(data);
		//prepareToSend["data"] = `"data":{"__schema":{"types":[{"name":"__DirectiveLocation","description":"","kind":"<introspection.TypeKind Value>"},{"kind":"<introspection.TypeKind Value>","name":"String","description":"The Stringscalar type represents textual data, represented as UTF-8 character sequences. The String type is most often used by GraphQL to represent free-form human-readable text."},{"kind":"<introspection.TypeKind Value>","name":"Boolean","description":"The Boolean scalar type represents true or false."},{"kind":"<introspection.TypeKind Value>","name":"__Schema","description":""},{"kind":"<introspection.TypeKind Value>","name":"__Directive","description":""},{"kind":"<introspection.TypeKind Value>","name":"__EnumValue","description":""},{"kind":"<introspection.TypeKind Value>","name":"Queries","description":""},{"kind":"<introspection.TypeKind Value>","name":"InputobEmp","description":""},{"kind":"<introspection.TypeKind Value>","name":"Int","description":"The Int scalar type represents non-fractional signed whole numeric values. Int can represent values between -(2^31) and 2^31 - 1."},{"kind":"<introspection.TypeKind Value>","name":"Float","description":"The Float scalar type represents signed double-precision fractional values as specified by [IEEE 754](http://en.wikipedia.org/wiki/IEEE_floating_point)."},{"kind":"<introspection.TypeKind Value>","name":"ID","description":"The ID scalar type represents a unique identifier, often used to refetch an object or as key for a cache. The ID type appears in a JSON response as a String; however, it is not intended to be human-readable. When expected as an input type, any string (such as \"4\") or integer (such as 4) input value will be accepted as an ID."},{"description":"","kind":"<introspection.TypeKind Value>","name":"__Type"},{"kind":"<introspection.TypeKind Value>","name":"Empleado","description":""},{"kind":"<introspection.TypeKind Value>","name":"Cargo","description":""},{"kind":"<introspection.TypeKind Value>","name":"__Field","description":""},{"kind":"<introspection.TypeKind Value>","name":"__InputValue","description":""},{"kind":"<introspection.TypeKind Value>","name":"__TypeKind","description":""},{"kind":"<introspection.TypeKind Value>","name":"Mutations","description":""}],"types2":[{"name":"Queries"},{"name":"__EnumValue"},{"name":"Float"},{"name":"ID"},{"name":"__Type"},{"name":"InputobEmp"},{"name":"Int"},{"name":"__InputValue"},{"name":"__TypeKind"},{"name":"Mutations"},{"name":"Empleado"},{"name":"Cargo"},{"name":"__Field"},{"name":"Boolean"},{"name":"__Schema"},{"name":"__Directive"},{"name":"__DirectiveLocation"},{"name":"String"}]},"pepe":{"types":[{"name":"Int"},{"name":"Float"},{"name":"ID"},{"name":"__Type"},{"name":"InputobEmp"},{"name":"__Field"},{"name":"__InputValue"},{"name":"__TypeKind"},{"name":"Mutations"},{"name":"Empleado"},{"name":"Cargo"},{"name":"String"},{"name":"Boolean"},{"name":"__Schema"},{"name":"__Directive"},{"name":"__DirectiveLocation"},{"name":"__EnumValue"},{"name":"Queries"}]}}`;
	}
	
	return prepareToSend
}
func(o *gql)setVariables(schema *ast.Schema, operation *ast.OperationDefinition,variables map[string]interface{}){
	//las operaciones tambien pueden tener directivas
	vars,err := validator.VariableValues(o.schema,operation,variables)
	//validar las variables con los scalar propios
	if err != nil{
		//responder con error
	}
	o.variables = vars;
}

func(o *gql) selectionSetParse(parse ast.SelectionSet, parent interface{}, parentProceced interface{}, typeName *string) ( Response ){
	
	//var prepareToSend Response
	prepareToSend := make(map[string]interface{},0);
	//ejecucion de las queries internas
	for _,selection :=range parse{
		rField := reflect.ValueOf(selection);
		switch rField.Type(){
		case reflect.TypeOf(&ast.Field{}):
			field := selection.(*ast.Field);
			prepareToSend = o.selectionParse(field, parent, parentProceced,typeName);
		case reflect.TypeOf(&ast.FragmentSpread{}):
			fragment := selection.(*ast.FragmentSpread);
			fragmentDef := fragment.Definition;
			for _,fragmentSelection := range fragmentDef.SelectionSet{
				field := fragmentSelection.(*ast.Field);
				prepareToSend = o.selectionParse(field, parent, parentProceced,typeName);
			}
		}
	}
	return prepareToSend;
}
func(o *gql) selectionParse(field *ast.Field, parent interface{}, parentProceced interface{}, typeName *string) (map[string]interface{}){
	fieldElem := field.Definition.Type.Elem;
	prepareToSend := make(map[string]interface{},0);
	var resolved resolvers.DataReturn;
	var resolvedProcesed resolvers.DataReturn;
	if field.SelectionSet != nil{
		namedType := field.Definition.Type.NamedType;
		fieldNames := o.getFieldNames(field.SelectionSet);
		if fieldElem != nil{
			namedType = fieldElem.NamedType;
		}
		if o.objectTypes[namedType] != nil{	
			args:= o.parseArguments(field.Arguments)
			directives := o.parseDirectives(field.Directives,namedType, field.Name);
			o.parseDirectives(field.Directives,namedType, field.Name);
			if typeName == nil{
				typeName = &namedType;
			}
			resolved = o.objectTypes[namedType].Resolver( field.Name, args, parent, directives , *typeName)
			resolvedProcesed = o.dataResponse(fieldNames, resolved);
		}
		rType :=  reflect.TypeOf(resolved);
		if rType != nil{
			rKind := rType.Kind();
			switch rKind{
				case reflect.Slice:
					var data []interface{};
					for i,value := range resolved.([]interface{}){
						responsed := o.selectionSetParse(field.SelectionSet,value, resolvedProcesed.([]interface{})[i], typeName);
						data = append(data,responsed);
					}
					if parentProceced != nil{
						prepareToSend = parentProceced.(map[string]interface{});
					}
					prepareToSend[field.Alias] = data;
				case reflect.Struct,reflect.Ptr:
					responsed := o.selectionSetParse(field.SelectionSet,resolved, resolvedProcesed, typeName);
					if parentProceced != nil{
						prepareToSend = parentProceced.(map[string]interface{});
					}
					prepareToSend[field.Alias] = responsed
				default:
					fmt.Println("aqui")
			}
		} else{
			if parentProceced != nil{
				prepareToSend = parentProceced.(map[string]interface{});
			}
			prepareToSend[field.Alias] = nil;
			//prepareToSend = parentProceced.(map[string]interface{});
		}
		
	} else{
		prepareToSend = parentProceced.(map[string]interface{});
	}
	return prepareToSend;
}
func(o *gql) parseDirectives(directives ast.DirectiveList, typeName string, fieldName string) (r resolvers.DirectiveList){
	r = make(map[string]interface{},0);
	for _,directive := range directives{
		args := make(map[string]interface{},0);
		for _,arg := range directive.Arguments{
			args[arg.Name] = arg;
		}
		var x resolvers.DataReturn;
		if o.directives[directive.Name] != nil{
			x = o.directives[directive.Name].Invoke(args,typeName,fieldName, o.schema.Directives[directive.Name]);
		}
		r[directive.Name] = x;
	}
	return r
}
func (o *gql) parseArguments(rawArgs ast.ArgumentList) map[string]interface{} {
	// en la definicion de los Kind hay que agregar la lista entera que proporciona la broma esa
	args := make(map[string]interface{})
	if len(rawArgs) > 0{
		for _,vArgs := range rawArgs{
			if len(vArgs.Value.Children) > 0 {
				args[vArgs.Name] = o.parseArgChildren(vArgs.Value.Children)
			}else{
				switch(vArgs.Value.Definition.Kind){
				case "OBJECT":	
					args[vArgs.Name] = vArgs.Value.Raw
				case "INTERFACE":
					args[vArgs.Name] = vArgs.Value.Raw
				case "UNION":
					args[vArgs.Name] = vArgs.Value.Raw
				case "ENUM":
					args[vArgs.Name] = vArgs.Value.Raw
				case "INPUT_OBJECT":
					args[vArgs.Name] = o.variables[vArgs.Value.Raw]
				case "SCALAR":
					switch(vArgs.Value.Definition.Name){
					case "String":
						args[vArgs.Name] = vArgs.Value.Raw
					case "Int":
						args[vArgs.Name] = vArgs.Value.Raw
					case "Float":
						args[vArgs.Name] = vArgs.Value.Raw
					case "Boolean":
						var data bool;
						if vArgs.Value.Raw == "true"{
							data = true;
						}
						args[vArgs.Name] = data
					case "ID":
						args[vArgs.Name] = vArgs.Value.Raw
					default:
						args[vArgs.Name] = vArgs.Value.Raw
					}
				default:
					args[vArgs.Name] = vArgs.Value.Raw
				}
			}
		}
	}
	return args;
}
func (o *gql) parseArgChildren(rawArgs ast.ChildValueList) map[string]interface{} {
	args := make(map[string]interface{})
	if len(rawArgs) > 0{
		for _,vArgs := range rawArgs{
			if len(vArgs.Value.Children) > 0 {
				args[vArgs.Name] = o.parseArgChildren(vArgs.Value.Children)
			}else{
				switch(vArgs.Value.Definition.Kind){
				case "OBJECT":	
					args[vArgs.Name] = vArgs.Value.Raw
				case "INTERFACE":
					args[vArgs.Name] = vArgs.Value.Raw
				case "UNION":
					args[vArgs.Name] = vArgs.Value.Raw
				case "ENUM":
					args[vArgs.Name] = vArgs.Value.Raw
				case "INPUT_OBJECT":
					args[vArgs.Name] = o.variables[vArgs.Value.Raw]
				case "SCALAR":
					switch(vArgs.Value.Definition.Name){
					case "String":
						args[vArgs.Name] = vArgs.Value.Raw
					case "Int":
						args[vArgs.Name] = vArgs.Value.Raw
					case "Float":
						args[vArgs.Name] = vArgs.Value.Raw
					case "Boolean":
						var data bool;
						if args[vArgs.Name] == "true"{
							data = true;
						}
						args[vArgs.Name] = data
					case "ID":
						args[vArgs.Name] = vArgs.Value.Raw
					default:
						args[vArgs.Name] = vArgs.Value.Raw
					}
				default:
					args[vArgs.Name] = vArgs.Value.Raw
				}
			}
		}
	}
	return args;
}
func(o *gql) getFieldNames(parse ast.SelectionSet) map[string]interface{}{
	fields := make(map[string]interface{});
	//debo anadir la consulta al 
	for _,selection :=range parse{
		rValue := reflect.ValueOf(selection);
		switch(rValue.Type()){
		case reflect.TypeOf(&ast.Field{}):
			field := selection.(*ast.Field);
			if field.Directives !=nil {
				//o.setDirectives(field.Name, field.Directives);
			}
			fields[field.Name] = field.Alias;
		case reflect.TypeOf(&ast.FragmentSpread{}):
			fragment := selection.(*ast.FragmentSpread);
			fragmentDef := fragment.Definition;
			for _,fragmentSelection := range fragmentDef.SelectionSet{
				field := fragmentSelection.(*ast.Field);
				fields[field.Name] = field.Alias;
			}

		case reflect.TypeOf(ast.InlineFragment{}):
		}
		
		
	}

	return fields;
}
