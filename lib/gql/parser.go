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
	maps "golang.org/x/exp/maps"
)

type Response interface{}

func (o *gql) response(request HttpRequest) (response HttpResponse) {
	document, err := gqlparser.LoadQuery(o.schema, request.Query)

	if err != nil {
		fmt.Println(err.Error())
	}
	if document != nil {
		parse := document.Operations
		if strings.Trim(request.OperationName, " ") != "" {
			for _, operation := range document.Operations {
				if operation.Name == request.OperationName {
					parse = ast.OperationList{operation}
					break
				}
			}
		}
		rx := o.operationParse(parse, request.Variables)
		response.Data = fmt.Sprintf("%v", rx["data"])
	}
	return response
}
func (o *gql) WebsocketResponse(request HttpRequest, socketId string, requestID RequestID, mt int) {
	document, err := gqlparser.LoadQuery(o.schema, request.Query)
	if err != nil {
		fmt.Println(err.Error())
	}
	if document != nil {
		parse := document.Operations
		if strings.Trim(request.OperationName, " ") != "" {
			for _, operation := range document.Operations {
				if operation.Name == request.OperationName {
					parse = ast.OperationList{operation}
					break
				}
			}
		}
		for _, operation := range parse {
			go o.websocketOperationParse(operation, request.Variables, socketId, requestID, mt)
		}
	}
}
func (o *gql) websocketOperationParse(operation *ast.OperationDefinition, variables map[string]interface{}, socketId string, requestID RequestID, mt int) {
	uuid := uuid.New().String()
	operationID := OperationID(operation.SelectionSet[0].(*ast.Field).Name)
	PubSub.createExcecuteEvent(EventID(uuid), operationID, SocketID(socketId), requestID, mt)
	response := &HttpResponse{}
	while := true
	for while {
		listen := PubSub.listenExcecuteEvent(operationID, EventID(uuid))
		rType := reflect.ValueOf(listen).Type()

		switch rType {
		case reflect.TypeOf(&SubscriptionClose{}):
			while = false
			break
		default:
			vars := o.setVariables(o.schema, operation, variables)
			var dataReturn resolvers.DataReturn
			data := make(map[string]interface{}, 0)
			isSubscriptionResponse := false
			switch operation.Operation {
			case ast.Subscription:
				data["data"], isSubscriptionResponse = o.selectionSetParse(string(operation.Operation), operation.SelectionSet, dataReturn, dataReturn, nil, 0, listen, vars)
			}
			if isSubscriptionResponse {
				response.Data = fmt.Sprintf("%v", o.jsonResponse(data))
				o.WriteWebsocketMessage(mt, socketId, requestID, response)
			}
		}
	}
}
func (o *gql) operationParse(parse ast.OperationList, variables map[string]interface{}) map[string]interface{} {
	prepareToSend := make(map[string]interface{}, 0)
	for _, operation := range parse {
		vars := o.setVariables(o.schema, operation, variables)
		var dataReturn resolvers.DataReturn
		data := make(map[string]interface{}, 0)
		switch operation.Operation {
		case ast.Query, ast.Mutation:
			data["data"], _ = o.selectionSetParse(string(operation.Operation), operation.SelectionSet, dataReturn, dataReturn, nil, 0, nil, vars)
		}
		prepareToSend["data"] = o.jsonResponse(data)
	}
	return prepareToSend
}
func (o *gql) setVariables(schema *ast.Schema, operation *ast.OperationDefinition, variables map[string]interface{}) (r map[string]any) {
	//las operaciones tambien pueden tener directivas
	vars, err := validator.VariableValues(o.schema, operation, variables)
	//validar las variables con los scalar propios
	if err != nil {
		fmt.Println(err)
	}
	return vars
}

func (o *gql) selectionSetParse(operation string, parse ast.SelectionSet, parent interface{}, parentProceced interface{}, typeName *string, start int, subscriptionValue interface{}, vars map[string]any) (Response, bool) {

	//var prepareToSend Response
	isSubscriptionResponse := false
	prepareToSend := make(map[string]interface{}, 0)
	send := make(map[string]interface{}, 0)
	//ejecucion de las queries internas
	for _, selection := range parse {
		rField := reflect.ValueOf(selection)
		switch rField.Type() {
		case reflect.TypeOf(&ast.Field{}):
			field := selection.(*ast.Field)
			prepareToSend, isSubscriptionResponse = o.selectionParse(operation, field, parent, parentProceced, typeName, start, subscriptionValue, vars)
		case reflect.TypeOf(&ast.FragmentSpread{}):
			fragment := selection.(*ast.FragmentSpread)
			fragmentDef := fragment.Definition
			for _, fragmentSelection := range fragmentDef.SelectionSet {
				field := fragmentSelection.(*ast.Field)
				prepareToSend, isSubscriptionResponse = o.selectionParse(operation, field, parent, parentProceced, typeName, start, subscriptionValue, vars)
			}
		case reflect.TypeOf(&ast.InlineFragment{}):
			fragment := selection.(*ast.InlineFragment)
			for _, fragmentSelection := range fragment.SelectionSet {
				field := fragmentSelection.(*ast.Field)
				prepareToSend, isSubscriptionResponse = o.selectionParse(operation, field, parent, parentProceced, typeName, start, subscriptionValue, vars)
			}
		}
		if start == 0 {
			maps.Copy(send, prepareToSend)
			prepareToSend = send
		}
	}
	return prepareToSend, isSubscriptionResponse
}
func (o *gql) selectionParse(operation string, field *ast.Field, parent interface{}, parentProceced interface{}, typeName *string, start int, subscriptionValue interface{}, vars map[string]any) (map[string]interface{}, bool) {
	fieldElem := field.Definition.Type.Elem
	isSubscriptionResponse := false
	prepareToSend := make(map[string]interface{}, 0)
	var resolved resolvers.DataReturn
	var resolvedProcesed resolvers.DataReturn
	if field.SelectionSet != nil {
		namedType := field.Definition.Type.NamedType
		fieldNames, typeCondition := o.getFieldNames(field.SelectionSet)
		if fieldElem != nil {
			namedType = fieldElem.NamedType
		}
		if typeCondition {
			//namedType = typeCondition
		}
		if o.objectTypes[namedType] != nil {
			args := o.parseArguments(field.Arguments, field.Definition.Arguments, vars)
			directives := o.parseDirectives(field.Directives, namedType, field.Name, vars)
			resolverInfo := resolvers.ResolverInfo{
				Operation:         operation,
				Args:              args,
				Resolver:          field.Name,
				Parent:            parent,
				Directives:        directives,
				TypeName:          namedType,
				ParentTypeName:    typeName,
				SubscriptionValue: subscriptionValue,
			}
			if operation == "subscription" && start == 0 {
				if ok := o.objectTypes[namedType].Subscribe(resolverInfo); ok {
					resolved, _ = o.objectTypes[namedType].Resolver(resolverInfo)
					typeName = &namedType
					resolvedProcesed = o.dataResponse(fieldNames, resolved, namedType)
					isSubscriptionResponse = true
				}
			} else {
				resolved, _ = o.objectTypes[namedType].Resolver(resolverInfo)
				typeName = &namedType
				resolvedProcesed = o.dataResponse(fieldNames, resolved, namedType)
			}
		}
		rType := reflect.TypeOf(resolved)
		if rType != nil {
			rKind := rType.Kind()
			switch rKind {
			case reflect.Slice:
				var data []interface{}
				rValue := reflect.ValueOf(resolved)
				for i := 0; i < rValue.Len(); i++ {
					responsed, _ := o.selectionSetParse(operation, field.SelectionSet, rValue.Index(i).Interface(), resolvedProcesed.([]interface{})[i], typeName, 1, subscriptionValue, vars)
					data = append(data, responsed)
				}
				if parentProceced != nil {
					prepareToSend = parentProceced.(map[string]interface{})
				}
				prepareToSend[field.Alias] = data
			case reflect.Struct, reflect.Ptr:
				responsed, _ := o.selectionSetParse(operation, field.SelectionSet, resolved, resolvedProcesed, typeName, 1, subscriptionValue, vars)
				if parentProceced != nil {
					prepareToSend = parentProceced.(map[string]interface{})
				}
				prepareToSend[field.Alias] = responsed
			}
		} else {
			if parentProceced != nil {
				prepareToSend = parentProceced.(map[string]interface{})
			}
			prepareToSend[field.Alias] = nil
			//prepareToSend = parentProceced.(map[string]interface{});
		}

	} else {
		prepareToSend = parentProceced.(map[string]interface{})

	}
	return prepareToSend, isSubscriptionResponse
}
func (o *gql) parseDirectives(directives ast.DirectiveList, typeName string, fieldName string, vars map[string]any) (r resolvers.DirectiveList) {
	r = make(map[string]interface{}, 0)
	for _, directive := range directives {
		args := o.parseArguments(directive.Arguments, directive.Definition.Arguments, vars)
		var x resolvers.DataReturn
		if o.directives[directive.Name] != nil {
			x, _ = o.directives[directive.Name].Invoke(args, typeName, fieldName)
		}
		r[directive.Name] = x
	}
	return r
}
func (o *gql) getFieldNames(parse ast.SelectionSet) (fields map[string]interface{}, typeCondition bool) {
	fields = make(map[string]interface{})
	//debo anadir la consulta al
	for _, selection := range parse {
		rValue := reflect.ValueOf(selection)
		switch rValue.Type() {
		case reflect.TypeOf(&ast.Field{}):
			field := selection.(*ast.Field)
			if field.Directives != nil {
				//o.setDirectives(field.Name, field.Directives);
			}
			fields[field.Name] = field.Alias
		case reflect.TypeOf(&ast.FragmentSpread{}):
			fragment := selection.(*ast.FragmentSpread)
			fragmentDef := fragment.Definition
			for _, fragmentSelection := range fragmentDef.SelectionSet {
				field := fragmentSelection.(*ast.Field)
				fields[field.Name] = field.Alias
			}
		case reflect.TypeOf(&ast.InlineFragment{}):

			fragment := selection.(*ast.InlineFragment)
			fields[fragment.TypeCondition] = fragment.TypeCondition
			typeCondition = true
			/*
				for _, fragmentSelection := range fragment.SelectionSet {
					field := fragmentSelection.(*ast.Field)
					fields[field.Name] = field.Alias
				}
			*/
		}

	}

	return
}
