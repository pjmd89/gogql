package definitionError

type errorLevel int

const (
	LEVEL_WARNING errorLevel = iota
	LEVEL_FATAL
)

type ErrorDescriptor struct {
	Message string
	Code    string
	Level   errorLevel
}
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
	Code       string
	Message    string             `json:"message"`
	Locations  []GQLErrorLocation `json:"locations,omitempty"`
	Path       []interface{}      `json:"path,omitempty"`
	Extensions ExtensionError     `json:"extensions,omitempty"`
}

type GQLError interface {
	Error() ErrorStruct
}
type ErrorList []GQLError
type Warning struct {
	GQLError
	ErrorStruct
}
type Fatal struct {
	GQLError
	ErrorStruct
}
