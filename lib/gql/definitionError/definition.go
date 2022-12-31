package definitionError

type ErrorLocation struct {
	line   int
	column int
}
type ExtensionError map[string]interface{}

type GQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
type ErrorStruct struct {
	Message    string             `json:"message"`
	Locations  []GQLErrorLocation `json:"locations,omitempty"`
	Path       []interface{}      `json:"path,omitempty"`
	Code       string             `json:"code,omitempty"`
	Extensions ExtensionError     `json:"extensions,omitempty"`
}

type GQLError interface {
	Error() ErrorStruct
}
type ErrorList []GQLError
type Fatal struct {
	GQLError
	ErrorStruct
}
type Warning struct {
	GQLError
	ErrorStruct
}
