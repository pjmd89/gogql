package definitionError

func NewError(err ErrorDescriptor, extensions ExtensionError) (r GQLError) {
	switch err.Level {
	case LEVEL_FATAL:
		r = &Fatal{
			ErrorStruct: ErrorStruct{
				Message:    err.Message,
				Code:       err.Code,
				Extensions: extensions,
			},
		}
	default:
		r = &Warning{
			ErrorStruct: ErrorStruct{
				Message:    err.Message,
				Code:       err.Code,
				Extensions: extensions,
			},
		}
	}
	return
}
func NewWarning(message string, extensions ExtensionError) *Warning {
	r := &Warning{
		ErrorStruct: ErrorStruct{Message: message, Extensions: extensions, Code: "000"},
	}
	return r
}
func NewFatal(message string, extensions ExtensionError) *Fatal {
	r := &Fatal{
		ErrorStruct: ErrorStruct{Message: message, Extensions: extensions, Code: "000"},
	}
	return r
}
func (o *Warning) Error() (r ErrorStruct) {
	r = o.ErrorStruct
	return
}
func (o *Fatal) Error() (r ErrorStruct) {
	r = o.ErrorStruct
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
