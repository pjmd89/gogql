package definitionError

type ErrorLocation struct {
	line   int
	column int
}
type ExtensionError map[string]interface{}
type errorStruct struct {
	message    string
	location   []ErrorLocation
	path       []interface{}
	extensions ExtensionError
}
type GQLError interface {
	GetError() errorStruct
}
type ErrorList []GQLError
type Fatal struct {
	GQLError
	errorStruct
}
type Warning struct {
	GQLError
	errorStruct
}
