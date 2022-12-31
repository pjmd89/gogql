package definitionError

func NewWarning(message string, extensions ExtensionError) *Warning {
	r := &Warning{
		ErrorStruct: ErrorStruct{Message: message, Extensions: extensions},
	}
	return r
}

func NewFatal(message string, extensions ExtensionError) *Fatal {
	r := &Fatal{
		ErrorStruct: ErrorStruct{Message: message, Extensions: extensions},
	}
	return r
}

func (o *Warning) Error() (r ErrorStruct) {
	r = ErrorStruct{o.Message, o.Locations, o.Path, o.Code, o.Extensions}
	return
}
func (o *Fatal) Error() (r ErrorStruct) {
	r = ErrorStruct{o.Message, o.Locations, o.Path, o.Code, o.Extensions}
	return
}
func (o ErrorList) GetErrors() (r []ErrorStruct) {
	if len(o) > 0 {
		r = make([]ErrorStruct, 0)
		for _, v := range o {
			r = append(r, v.Error())
		}
	}
	return
}
