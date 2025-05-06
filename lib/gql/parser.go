package gql

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/resolvers"
	"github.com/pjmd89/gqlparser/v2"
	"github.com/pjmd89/gqlparser/v2/ast"
	"github.com/pjmd89/gqlparser/v2/validator"
	maps "golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Response interface{}

var introspectionType = []string{
	"__Schema",
	"__Type",
	"__Field",
	"__EnumValue",
	"__InputValue",
	"__Directive",
}
var introspectionResolver = []string{
	"__schema",
	"__type",
}

func (o *Gql) response(request HttpRequest, sessionID string) (response HttpResponse) {
	document, err := gqlparser.LoadQuery(o.schema, request.Query)
	if err != nil {
		log.Println(err.Error())
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
		response = o.operationParse(parse, request.Variables, sessionID)
		//response.Data = fmt.Sprintf("%v", rx["data"])
	}
	return response
}
func (o *Gql) WebsocketResponse(request HttpRequest, socketId string, requestID RequestID, mt int, sessionID string) {
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
			go o.websocketOperationParse(operation, request.Variables, socketId, requestID, mt, sessionID)
		}
	}
}
func (o *Gql) websocketOperationParse(operation *ast.OperationDefinition, variables map[string]interface{}, socketId string, requestID RequestID, mt int, sessionID string) {
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
		default:
			vars, _ := o.setVariables(operation, variables)
			var dataReturn resolvers.DataReturn
			data := make(map[string]interface{}, 0)
			isSubscriptionResponse := false
			switch operation.Operation {
			case ast.Subscription:
				errList := make(definitionError.ErrorList, 0)
				data["data"], isSubscriptionResponse = o.selectionSetParse(string(operation.Operation), operation.SelectionSet, dataReturn, dataReturn, nil, 0, listen, vars, sessionID, &errList)
			}
			if isSubscriptionResponse {
				response.Data = fmt.Sprintf("%v", o.jsonResponse(data))
				o.WriteWebsocketMessage(mt, socketId, requestID, response)
			}
		}
	}
}
func (o *Gql) operationParse(parse ast.OperationList, variables map[string]interface{}, sessionID string) HttpResponse {
	response := HttpResponse{}
	errList := make(definitionError.ErrorList, 0)
	for _, operation := range parse {
		var dataReturn resolvers.DataReturn
		var data Response
		vars, setErr := o.setVariables(operation, variables)
		if setErr != nil {
			log.Println(setErr.Error())
			break
		}
		switch operation.Operation {
		case ast.Query, ast.Mutation:
			data, _ = o.selectionSetParse(string(operation.Operation), operation.SelectionSet, dataReturn, dataReturn, nil, 0, nil, vars, sessionID, &errList)
		}
		response.Data = fmt.Sprintf("%v", o.jsonResponse(data))
	}
	if len(errList) > 0 {
		errString, _ := json.Marshal(errList.GetErrors())
		response.Errors = string(errString)
	}
	return response
}
func (o *Gql) setVariables(operation *ast.OperationDefinition, variables map[string]interface{}) (r map[string]any, err error) {
	//las operaciones tambien pueden tener directivas
	//validar las variables con los scalar propios
	r, varsErr := validator.VariableValues(o.schema, operation, variables)
	if varsErr != nil {
		err = varsErr
	}
	return
}

func (o *Gql) selectionSetParse(operation string, parse ast.SelectionSet, parent interface{}, parentProceced interface{}, typeName *string, start int, subscriptionValue interface{}, vars map[string]any, sessionID string, errList *definitionError.ErrorList) (Response, bool) {

	//var prepareToSend Response
	isSubscriptionResponse := false
	prepareToSend := make(map[string]interface{}, 0)
	send := make(map[string]interface{}, 0)
	//ejecucion de las queries internas
	for _, selection := range parse {
		var stop, fieldTypeIsUnion bool
		rField := reflect.ValueOf(selection)
		switch rField.Type() {
		case reflect.TypeOf(&ast.Field{}):
			field := selection.(*ast.Field)
			fieldTypeIsUnion = o.typeIsUnion(field)
			prepareToSend, isSubscriptionResponse, stop = o.selectionParse(operation, field, parent, parentProceced, typeName, start, subscriptionValue, vars, sessionID, errList, fieldTypeIsUnion)
		case reflect.TypeOf(&ast.FragmentSpread{}):
			fragment := selection.(*ast.FragmentSpread)
			fragmentDef := fragment.Definition
			for _, fragmentSelection := range fragmentDef.SelectionSet {
				field := fragmentSelection.(*ast.Field)
				fieldTypeIsUnion = o.typeIsUnion(field)
				prepareToSend, isSubscriptionResponse, stop = o.selectionParse(operation, field, parent, parentProceced, typeName, start, subscriptionValue, vars, sessionID, errList, fieldTypeIsUnion)
			}
		case reflect.TypeOf(&ast.InlineFragment{}):
			fragment := selection.(*ast.InlineFragment)
			for _, fragmentSelection := range fragment.SelectionSet {
				field := fragmentSelection.(*ast.Field)
				fieldTypeIsUnion = o.typeIsUnion(field)
				prepareToSend, isSubscriptionResponse, stop = o.selectionParse(operation, field, parent, parentProceced, typeName, start, subscriptionValue, vars, sessionID, errList, fieldTypeIsUnion)
			}
		}
		if start == 0 {
			maps.Copy(send, prepareToSend)
			prepareToSend = send
		}
		if stop {
			break
		}
	}
	return prepareToSend, isSubscriptionResponse
}

func (o *Gql) typeIsUnion(field *ast.Field) (isUnion bool) {
	namedType := field.Definition.Type.NamedType
	if field.Definition.Type.Elem != nil {
		namedType = field.Definition.Type.Elem.NamedType
	}

	if val, ok := o.schema.Types[namedType]; ok && val.Kind == "UNION" {
		isUnion = true
	}

	return
}
func (o *Gql) selectionParse(operation string, field *ast.Field, parent interface{}, parentProceced interface{}, typeName *string, start int, subscriptionValue interface{}, vars map[string]any, sessionID string, errList *definitionError.ErrorList, isUnion bool) (map[string]interface{}, bool, bool) {
	fieldElem := field.Definition.Type.Elem
	isArr := field.Definition.Type.Elem != nil
	isSubscriptionResponse := false
	canResolve := true
	prepareToSend := make(map[string]interface{}, 0)
	var resolved, resolvedProcesed resolvers.DataReturn
	var err definitionError.GQLError
	namedType := field.Definition.Type.NamedType

	if field.SelectionSet != nil {
		fieldNames, _, fieldGroup := o.getFieldNames(field.SelectionSet)
		if fieldElem != nil {
			namedType = fieldElem.NamedType
		}
		var authenticatedError definitionError.GQLError
		if !slices.Contains(introspectionType, namedType) {
			authInfo := AuthorizateInfo{operation, TypeName(field.ObjectDefinition.Name), TypeName(namedType), ResolverName(field.Name), sessionID}
			authenticatedError = o.OnAuthorizate(authInfo)
		}
		if authenticatedError == nil && slices.Contains(introspectionType, namedType) && slices.Contains(introspectionResolver, field.Name) {
			authenticatedError = o.OnIntrospection()
		}

		if authenticatedError != nil {
			resolvedProcesed = o.dataResponse(fieldNames, nil, namedType)
			*errList = append(*errList, authenticatedError)
			if reflect.TypeOf(authenticatedError) == reflect.TypeOf(&definitionError.Fatal{}) {
				return nil, isSubscriptionResponse, true
			}
		}

		if o.objectTypes[namedType] != nil && authenticatedError == nil {
			args, vError := o.parseArguments(field.Arguments, field.Definition.Arguments, vars)
			var argumentError definitionError.GQLError
			if vError != nil {
				*errList = append(*errList, vError)
				argumentError = vError
			}
			directives, vError := o.parseDirectives(field.Directives, namedType, field.Name, vars)
			if vError != nil {
				*errList = append(*errList, vError)
				argumentError = vError
			}

			if argumentError != nil {
				return nil, isSubscriptionResponse, true
			}
			resolverInfo := resolvers.ResolverInfo{
				Operation:         operation,
				Args:              args,
				Resolver:          field.Name,
				Parent:            parent,
				Directives:        directives,
				TypeName:          namedType,
				ParentTypeName:    typeName,
				SubscriptionValue: subscriptionValue,
				SessionID:         sessionID,
			}
			if operation == "subscription" && start == 0 {
				canResolve = o.objectTypes[namedType].Subscribe(resolverInfo)
				isSubscriptionResponse = canResolve
			}
			if canResolve {
				resolved, resolvedProcesed, err = o.resolve(fieldNames, fieldGroup, namedType, resolverInfo, isArr, isUnion)
				if err != nil {
					*errList = append(*errList, err)
				}
				if err != nil && err.ErrorLevel() == definitionError.LEVEL_FATAL {
					return nil, isSubscriptionResponse, true
				}
				typeName = &namedType
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
					responsed, _ := o.selectionSetParse(operation, field.SelectionSet, rValue.Index(i).Interface(), resolvedProcesed.([]interface{})[i], typeName, 1, subscriptionValue, vars, sessionID, errList)
					data = append(data, responsed)
				}
				if parentProceced != nil {
					prepareToSend = parentProceced.(map[string]interface{})
				}
				prepareToSend[field.Alias] = data
			case reflect.Struct, reflect.Ptr:
				responsed, _ := o.selectionSetParse(operation, field.SelectionSet, resolved, resolvedProcesed, typeName, 1, subscriptionValue, vars, sessionID, errList)
				if parentProceced != nil {
					prepareToSend = parentProceced.(map[string]interface{})
				}
				prepareToSend[field.Alias] = responsed
			}
		} else {
			if parentProceced != nil {
				prepareToSend = parentProceced.(map[string]interface{})
			}
			/*if _, ok := fieldGroup[unionType.(string)][varName]; !ok && varName != "__typename"{
				break
			}
			*/
			prepareToSend[field.Alias] = nil
			//prepareToSend = parentProceced.(map[string]interface{});
		}

	} else {
		if parentProceced != nil {
			prepareToSend = parentProceced.(map[string]interface{})
		} else {
			if field.Name == "__typename" {
				prepareToSend[field.Alias] = field.ObjectDefinition.Name
			}
		}
	}
	return prepareToSend, isSubscriptionResponse, false
}

func (o *Gql) resolve(fieldnames map[string]string, fieldGroup map[string]map[string]string, namedType string, resolverInfo resolvers.ResolverInfo, isArr, isUnion bool) (rawResolved, parsedResolved resolvers.DataReturn, err definitionError.GQLError) {
	if !isUnion {
		dataReturn, resolvErr := o.objectTypes[namedType].Resolver(resolverInfo)
		if resolvErr != nil && resolvErr.ErrorLevel() == definitionError.LEVEL_FATAL {
			err = resolvErr
			return
		}
		return dataReturn, o.dataResponse(fieldnames, dataReturn, namedType), resolvErr
	}

	unionData, err := o.objectTypes[namedType].Resolver(resolverInfo)
	if err != nil && err.ErrorLevel() == definitionError.LEVEL_FATAL {
		return
	}

	if unionData == nil || reflect.ValueOf(unionData).Type() != reflect.TypeOf(map[string]any{}) {
		log.Printf("return data of union %s resolver must be of type map[string]any\n", namedType)
		return
	}

	var resolversDataReturn, processedResolved []any
	for typeResolver, data := range unionData.(map[string]any) {
		if index := slices.Index(o.schema.Types[namedType].Types, typeResolver); index == -1 {
			log.Printf("map key returned by union %s resolver must be in the union types\n", namedType)
			continue
		}

		if _, ok := fieldGroup[typeResolver]; !ok || o.objectTypes[typeResolver] == nil {
			continue
		}

		resolverInfo.Parent = data
		dataReturn, resolvErr := o.objectTypes[typeResolver].Resolver(resolverInfo)
		if resolvErr != nil && resolvErr.ErrorLevel() == definitionError.LEVEL_FATAL {
			err = resolvErr
			return
		}

		if dataReturn != nil {
			switch reflect.TypeOf(dataReturn).Kind() {
			case reflect.Slice, reflect.Array:
				r := o.dataResponse(fieldGroup[typeResolver], dataReturn, typeResolver)
				processedResolved = append(processedResolved, r.([]any)...)
				rValueArray := reflect.ValueOf(dataReturn)
				for i := 0; i < rValueArray.Len(); i++ {
					resolversDataReturn = append(resolversDataReturn, rValueArray.Index(i).Interface())
				}
			default:
				resolversDataReturn = append(resolversDataReturn, dataReturn)
				processedResolved = append(processedResolved, o.dataResponse(fieldGroup[typeResolver], dataReturn, typeResolver))
			}
		}
	}

	if !isArr && len(resolversDataReturn) > 0 {
		rawResolved = resolversDataReturn[0]
		parsedResolved = processedResolved[0]
	}

	if isArr {
		rawResolved = resolversDataReturn
		parsedResolved = processedResolved
	}

	return
}
func (o *Gql) parseDirectives(directives ast.DirectiveList, typeName string, fieldName string, vars map[string]any) (r resolvers.DirectiveList, err definitionError.GQLError) {
	r = make(map[string]interface{}, 0)
	for _, directive := range directives {
		args, vError := o.parseArguments(directive.Arguments, directive.Definition.Arguments, vars)
		if vError != nil {
			err = vError
			return
		}
		var x resolvers.DataReturn
		if o.directives[directive.Name] != nil {
			x, _ = o.directives[directive.Name].Invoke(args, typeName, fieldName)
		}
		r[directive.Name] = x
	}
	return
}

func (o *Gql) getFieldNames(parse ast.SelectionSet) (fields map[string]string, typeCondition bool, fieldsGroup map[string]map[string]string) {
	fields = make(map[string]string, 0)
	fieldsGroup = map[string]map[string]string{}
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
				if _, ok := fieldsGroup[fragmentDef.TypeCondition]; !ok {
					fieldsGroup[fragmentDef.TypeCondition] = make(map[string]string)
				}
				fieldsGroup[fragmentDef.TypeCondition][field.Name] = field.Alias
			}
		case reflect.TypeOf(&ast.InlineFragment{}):
			fragment := selection.(*ast.InlineFragment)
			//fields[fragment.TypeCondition] = fragment.TypeCondition
			//typeCondition = true
			for _, fragmentSelection := range fragment.SelectionSet {
				field := fragmentSelection.(*ast.Field)
				fields[field.Name] = field.Alias
				if _, ok := fieldsGroup[fragment.TypeCondition]; !ok {
					fieldsGroup[fragment.TypeCondition] = make(map[string]string)
				}
				fieldsGroup[fragment.TypeCondition][field.Name] = field.Alias
			}
		}
	}
	return
}
