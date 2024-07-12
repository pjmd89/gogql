package definitionError

func NewError(err ErrorDescriptor, extensions ExtensionError) (r GQLError) {
	extensions = setExtension(extensions, err.Level, err.Code)
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
	extensions = setExtension(extensions, LEVEL_WARNING, "000")
	r := &Warning{
		ErrorStruct: ErrorStruct{Message: message, Extensions: extensions, Code: "000"},
	}
	return r
}
func NewFatal(message string, extensions ExtensionError) *Fatal {
	extensions = setExtension(extensions, LEVEL_FATAL, "000")
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
func setExtension(extensions ExtensionError, errLevel errorLevel, code string) (r ExtensionError) {
	levelName := map[errorLevel]string{LEVEL_FATAL: "fatal", LEVEL_WARNING: "warning"}
	if extensions == nil {
		extensions = map[string]any{}
	}
	if extensions["code"] == nil {
		extensions["code"] = code
	}
	extensions["level"] = levelName[errLevel]
	r = extensions
	return
}
