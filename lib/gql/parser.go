package gql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
	"github.com/pjmd89/gqlparser/v2"
	"github.com/pjmd89/gqlparser/v2/ast"
	"github.com/pjmd89/gqlparser/v2/validator"
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
func(o *gql) WebsocketResponse(request HttpRequest, socketId string, requestID RequestID, mt int){
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
		for _,operation :=range parse{
			go o.websocketOperationParse(operation,request.Variables, socketId, requestID, mt);
		}
	}
}
func(o *gql) websocketOperationParse(operation *ast.OperationDefinition, variables map[string]interface{}, socketId string, requestID RequestID ,mt int){
	uuid := uuid.New().String();
	operationID := OperationID(operation.SelectionSet[0].(*ast.Field).Name);
	PubSub.createExcecuteEvent(EventID(uuid),operationID,SocketID(socketId),requestID,mt);
	response := &HttpResponse{};
	while := true;
	for while{
		listen:= PubSub.listenExcecuteEvent(operationID,EventID(uuid));
		rType := reflect.ValueOf(listen).Type();
		
		switch rType{
		case reflect.TypeOf(&SubscriptionClose{}):
			while = false;
			break;
		default:
			o.setVariables(o.schema,operation,variables);
			var dataReturn resolvers.DataReturn;
			data := make(map[string]interface{},0);
			isSubscriptionResponse:= false;
			switch operation.Operation{
			case ast.Subscription:
				data["data"],isSubscriptionResponse = o.selectionSetParse(string(operation.Operation), operation.SelectionSet, dataReturn ,dataReturn,nil, 0, listen);
			}
			if isSubscriptionResponse{
				response.Data = fmt.Sprintf("%v",o.jsonResponse(data));
			o.WriteWebsocketMessage(mt,socketId,requestID, response);
			}
		}
	}
}
func(o *gql) operationParse(parse ast.OperationList, variables map[string]interface{}) (map[string]interface{}){
	prepareToSend := make(map[string]interface{},0);
	for _,operation :=range parse{
		o.setVariables(o.schema,operation,variables);
		var dataReturn resolvers.DataReturn;
		data := make(map[string]interface{},0);
		switch operation.Operation{
		case ast.Query, ast.Mutation:
			data["data"],_ = o.selectionSetParse(string(operation.Operation), operation.SelectionSet, dataReturn ,dataReturn,nil, 0, nil);
		}
		prepareToSend["data"] = o.jsonResponse(data);
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

func(o *gql) selectionSetParse(operation string, parse ast.SelectionSet, parent interface{}, parentProceced interface{}, typeName *string,start int, subscriptionValue interface{}) ( Response, bool ){
	
	//var prepareToSend Response
	isSubscriptionResponse := false;
	prepareToSend := make(map[string]interface{},0);
	//ejecucion de las queries internas
	for _,selection :=range parse{
		rField := reflect.ValueOf(selection);
		switch rField.Type(){
		case reflect.TypeOf(&ast.Field{}):
			field := selection.(*ast.Field);
			prepareToSend,isSubscriptionResponse = o.selectionParse(operation, field, parent, parentProceced,typeName, start, subscriptionValue );
		case reflect.TypeOf(&ast.FragmentSpread{}):
			fragment := selection.(*ast.FragmentSpread);
			fragmentDef := fragment.Definition;
			for _,fragmentSelection := range fragmentDef.SelectionSet{
				field := fragmentSelection.(*ast.Field);
				prepareToSend,isSubscriptionResponse = o.selectionParse(operation, field, parent, parentProceced,typeName, start, subscriptionValue );
			}
		}
	}
	return prepareToSend,isSubscriptionResponse;
}
func(o *gql) selectionParse( operation string, field *ast.Field, parent interface{}, parentProceced interface{}, typeName *string, start int, subscriptionValue interface{}) (map[string]interface{}, bool){
	fieldElem := field.Definition.Type.Elem;
	isSubscriptionResponse := false;
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
			args:= o.parseArguments(field.Arguments,field.Definition.Arguments);
			directives := o.parseDirectives(field.Directives,namedType, field.Name);
			resolverInfo := resolvers.ResolverInfo{
				Operation: operation,
				Args: args,
				Resolver: field.Name,
				Parent: parent,
				Directives: directives,
				TypeName: namedType,
				ParentTypeName: typeName,
				SubscriptionValue: subscriptionValue,
			};
			if operation == "subscription" && start == 0{
				if ok := o.objectTypes[namedType].Subscribe( resolverInfo ); ok{
					resolved = o.objectTypes[namedType].Resolver( resolverInfo );
					typeName = &namedType;
					resolvedProcesed = o.dataResponse(fieldNames, resolved);
					isSubscriptionResponse = true;
				}
			} else{
				resolved = o.objectTypes[namedType].Resolver( resolverInfo );
				typeName = &namedType;
				resolvedProcesed = o.dataResponse(fieldNames, resolved);
			}
		}
		rType :=  reflect.TypeOf(resolved);
		if rType != nil{
			rKind := rType.Kind();
			switch rKind{
				case reflect.Slice:
					var data []interface{};
					for i,value := range resolved.([]interface{}){
						responsed,_ := o.selectionSetParse(operation, field.SelectionSet,value, resolvedProcesed.([]interface{})[i], typeName, 1, subscriptionValue );
						data = append(data,responsed);
					}
					if parentProceced != nil{
						prepareToSend = parentProceced.(map[string]interface{});
					}
					prepareToSend[field.Alias] = data;
				case reflect.Struct,reflect.Ptr:
					responsed,_ := o.selectionSetParse(operation, field.SelectionSet,resolved, resolvedProcesed, typeName, 1, subscriptionValue );
					if parentProceced != nil{
						prepareToSend = parentProceced.(map[string]interface{});
					}
					prepareToSend[field.Alias] = responsed
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
	return prepareToSend,isSubscriptionResponse;
}
func(o *gql) parseDirectives(directives ast.DirectiveList, typeName string, fieldName string) (r resolvers.DirectiveList){
	r = make(map[string]interface{},0);
	for _,directive := range directives{
		args:= o.parseArguments(directive.Arguments, directive.Definition.Arguments);
		var x resolvers.DataReturn;
		if o.directives[directive.Name] != nil{
			x = o.directives[directive.Name].Invoke(args,typeName,fieldName);
		}
		r[directive.Name] = x;
	}
	return r
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
