package gql

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/structs"
	"github.com/pjmd89/gogql/lib/resolvers"
	"github.com/pjmd89/goutils/dbutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// *
func (o *Gql) dataResponse(fieldNames map[string]string, resolved interface{}, resolverName string) (r interface{}) {
	rType := reflect.TypeOf(resolved)
	if rType != nil {
		rKind := rType.Kind()
		switch rKind {
		case reflect.Slice, reflect.Array:
			switch rType {
			case reflect.TypeOf(primitive.ObjectID{}):
				rValue := reflect.ValueOf(resolved)
				id := rValue.Interface().(primitive.ObjectID).Hex()
				r = reflect.ValueOf(id).Interface()
			default:
				r = make([]interface{}, 0)
				rValue := reflect.ValueOf(resolved)
				for i := 0; i < rValue.Len(); i++ {
					value := o.dataResponse(fieldNames, rValue.Index(i).Interface(), resolverName)
					r = append(r.([]interface{}), o.dataResponse(fieldNames, value, resolverName))
				}
			}
		case reflect.Ptr:
			rValue := reflect.ValueOf(resolved)
			if !rValue.IsNil() {
				r = o.dataResponse(fieldNames, rValue.Elem().Interface(), resolverName)
			}

		case reflect.Struct:
			rValue := reflect.ValueOf(resolved)
			e := rValue.Type()

			r = make(map[string]interface{}, 0)
			if _, ok := fieldNames["__typename"]; ok {
				r.(map[string]interface{})["__typename"] = resolverName
			}

			for i := 0; i < e.NumField(); i++ {
				tagsContent := dbutils.GetTags(e.Field(i))
				varName := e.Field(i).Name
				if !tagsContent.IsOmit {
					if tagsContent.Name != "" {
						varName = tagsContent.Name
					}
					isNill := false
					switch e.Field(i).Type.Kind() {
					case reflect.Ptr:
						isNill = rValue.Field(i).IsNil()
					}
					if !isNill {
						data := o.dataResponse(make(map[string]string, 0), rValue.Field(i).Interface(), resolverName)
						if fieldNames[varName] != "" {
							y := e.Field(i).Type.Name()
							if o.scalars[y] != nil {
								resolvedData := resolvers.ScalarResolved{
									Value:        data,
									ResolverName: resolverName,
									Resolved:     structs.New(resolved),
								}
								data, _ = o.scalars[y].Assess(resolvedData)
							}
							r.(map[string]interface{})[fieldNames[varName]] = data
						}
					} else {
						if fieldNames[varName] != "" {
							r.(map[string]interface{})[fieldNames[varName]] = nil
						}
					}
				}
			}
		default:
			//aqui hay que hacer la evaluacion del scalar
			r = reflect.ValueOf(resolved).Interface()
		}
	}
	return r
}
func (o *Gql) getTags(e reflect.Type, i int) map[string]string {
	tagsContent := make(map[string]string, 0)
	varTag, _ := e.Field(i).Tag.Lookup("gql")
	splitTags := strings.Split(varTag, ",")
	if len(splitTags) > 0 {
		for _, value := range splitTags {
			splitValue := strings.Split(value, "=")
			if len(splitValue) > 1 {
				tagsContent[splitValue[0]] = splitValue[1]
			}
		}
	}
	return tagsContent
}

//*/

func (o *Gql) jsonResponse(data interface{}) interface{} {

	datax := o.prepareJson(data)
	return datax
}
func (o *Gql) prepareJson(data interface{}) interface{} {
	r := "null"
	rType := reflect.TypeOf(data)

	rValue := reflect.ValueOf(data)

	if rType != nil {
		rKind := rType.Kind()
		switch rKind {
		case reflect.Map:
			x := make([]string, 0)
			for i, value := range data.(map[string]interface{}) {
				datax := `"` + i + `":` + o.prepareJson(value).(string)
				x = append(x, datax)
			}
			r = ""
			if len(x) > 0 {
				r = "{" + strings.Join(x, ",") + "}"
			} else {
				r = "null"
			}

		case reflect.Struct:

		case reflect.Array | reflect.Slice:
			x := make([]string, 0)
			for _, value := range data.([]interface{}) {
				insert := true
				datax := o.prepareJson(value).(string)
				if reflect.ValueOf(value).Type().Kind() == reflect.Map && len(datax) == 0 {
					insert = false
				}
				if insert {
					x = append(x, datax)
				}

			}
			r = "[" + strings.Join(x, ",") + "]"
		default:
			switch rKind {
			case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint8,
				reflect.Uint32, reflect.Uint64:
				r = fmt.Sprint(rValue)
			default:
				///*
				scapeEnter := regexp.MustCompile(`[\n]`)
				scapeCarReturn := regexp.MustCompile(`[\r]`)
				scapeTab := regexp.MustCompile(`[\t]`)
				str := scapeEnter.ReplaceAllString(fmt.Sprint(rValue), `\n`)
				str = scapeCarReturn.ReplaceAllString(str, `\r`)
				str = scapeTab.ReplaceAllString(str, `\t`)
				r = `"` + strings.Replace(str, `"`, `\"`, -1) + `"`
				//*/
				//r = `"`+scape.ReplaceAllString(fmt.Sprint(rValue),`\n`)+`"`;
				//r = `"`+strings.Replace(rValue.String(),`"`,`\"`,-1)+`"`;
			}
			if rType.Name() == "TypeKind" {
				r = `"` + fmt.Sprint(rValue.Interface()) + `"`
			}

		}
	}
	return r
}
