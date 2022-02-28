package gql

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

//*
func(o *gql) dataResponse(fieldNames map[string]interface{}, resolved interface{}) (r interface{} ){
	rType :=  reflect.TypeOf(resolved);
	if rType != nil{
		rKind := rType.Kind();
		switch(rKind){
		case reflect.Slice, reflect.Array:
			r = make([]interface{},0);
			rValue := reflect.ValueOf(resolved);

			for i := 0; i < rValue.Len(); i++{
				value := o.dataResponse(fieldNames,rValue.Index(i).Interface());
				r = append(r.([]interface{}),o.dataResponse(fieldNames, value));
			}
		case reflect.Ptr:
			rValue := reflect.ValueOf(resolved);
			if !rValue.IsNil(){
				r = o.dataResponse(fieldNames,rValue.Elem().Interface());
			}
			
		case reflect.Struct:
			rValue := reflect.ValueOf(resolved);
			e := rValue.Type()
			
			r = make(map[string]interface{},0);
			for i:=0;i < e.NumField();i++{
				tagsContent := o.getTags(e,i);
				varName := e.Field(i).Name
				if tagsContent["name"] != ""{
					varName = tagsContent["name"];
				}
				isNill := false;
				switch e.Field(i).Type.Kind(){
					case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map:
						isNill = rValue.Field(i).IsNil();
				}
				if !isNill {
					data := o.dataResponse(make(map[string]interface{},0),rValue.Field(i).Interface());	
					if fieldNames[varName] != nil{
						r.(map[string]interface{})[fieldNames[varName].(string)] = data;
					}
				}else{
					if fieldNames[varName] != nil{
						r.(map[string]interface{})[fieldNames[varName].(string)] = nil;
					}
				}
			}
		default:
			//aqui hay que hacer la evaluacion del scalar 
			r = reflect.ValueOf(resolved).Interface();
		}
	}
	return r
}
func(o *gql) getTags(e reflect.Type, i int) map[string]string{
	tagsContent := make(map[string]string,0);
	varTag,_ := e.Field(i).Tag.Lookup("gql");
	splitTags := strings.Split(varTag,",");
	if(len(splitTags) > 0){
		for _,value := range splitTags{
			splitValue := strings.Split(value,"=");
			if(len(splitValue) > 1){
				tagsContent[splitValue[0]] = splitValue[1];
			}
		}
	}
	return tagsContent;
}
//*/

func(o *gql) jsonResponse( data interface{}) interface{}{
	
	datax := o.prepareJson(data);
	return datax;
}
func(o *gql) prepareJson(data interface{}) interface{}{
	r := "null";
	rType := reflect.TypeOf(data);
	
	rValue := reflect.ValueOf(data);

	if rType != nil{
		rKind := rType.Kind()
		switch(rKind){
			case reflect.Map:
				x := make([]string,0);
				for i,value :=range data.(map[string]interface{}){
					datax := `"`+i+`":`+o.prepareJson(value).(string);
					x = append(x,datax);
				}
				r = "";
				if len(x) > 0{
					r = "{"+strings.Join(x,",")+"}";
				}
				
			case reflect.Struct:
				
			case reflect.Array | reflect.Slice:
				x := make([]string,0);
				for _,value :=range data.([]interface{}){
					insert := true;
					datax := o.prepareJson(value).(string);
					if reflect.ValueOf(value).Type().Kind() == reflect.Map && len(datax) == 0{
						insert =false;
					}
					if insert{
						x = append(x,datax);
					}
					
				}
				r = "["+strings.Join(x,",")+"]"
			default:
				switch(rKind){
					case 	reflect.Bool,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,
							reflect.Float32,reflect.Float64, reflect.Uint, reflect.Uint8,
							reflect.Uint32, reflect.Uint64:
						r = fmt.Sprint(rValue)
					default:
						///*
						scapeEnter := regexp.MustCompile(`[\n]`);
						scapeCarReturn := regexp.MustCompile(`[\r]`);
						scapeTab := regexp.MustCompile(`[\t]`);
						str := scapeEnter.ReplaceAllString(fmt.Sprint(rValue),`\n`);
						str = scapeCarReturn.ReplaceAllString(str,`\r`);
						str = scapeTab.ReplaceAllString(str,`\t`);
						r = `"`+strings.Replace(str,`"`,`\"`,-1)+`"`;
						//*/
						//r = `"`+scape.ReplaceAllString(fmt.Sprint(rValue),`\n`)+`"`;
						//r = `"`+strings.Replace(rValue.String(),`"`,`\"`,-1)+`"`;
				}
				if rType.Name() == "TypeKind"{
					r = `"`+fmt.Sprint(rValue.Interface())+`"`;
				}
				
		}
	}
	return r;
}