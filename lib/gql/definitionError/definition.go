package definitionError
type ErrorLocation struct{
	line int
	column int
}
type ExtensionError map[string]interface{}
type errorStruct struct{
	message string
	location []ErrorLocation
	path []interface{}
    extensions ExtensionError
}
type Error interface{
	GetError() errorStruct
}
type ErrorList [] Error
type Fatal struct{
	Error
	errorStruct
}
type Warning struct{
	Error
	errorStruct
}
