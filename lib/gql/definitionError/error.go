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
	return r
}
