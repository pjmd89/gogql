package definitionError

func NewWarning(message string, extensions ExtensionError) *Warning {
	r := &Warning{
		errorStruct: errorStruct{message: message, extensions: extensions},
	}
	return r
}

func NewFatal(message string, extensions ExtensionError) *Fatal {
	r := &Fatal{
		errorStruct: errorStruct{message: message, extensions: extensions},
	}
	return r
}

func (o *Warning) Error() (r errorStruct) {
	return r
}
