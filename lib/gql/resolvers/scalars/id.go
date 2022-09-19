package scalars

import (
	"errors"
	"fmt"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type ID string

func NewIDScalar() (r resolvers.Scalar) {
	var scalar *ID
	r = scalar
	return
}
func (o *ID) Set(value interface{}) (r interface{}, err definitionError.GQLError) {
	s := fmt.Sprintf("%v", value)
	r = s
	return
}
func (o *ID) Assess(value interface{}) (val interface{}, err definitionError.GQLError) {
	var er error

	switch value.(type) {
	case string:
		val = value.(string)
	case int, int32, int64:
		val = value.(int64)
	case float32, float64:
		val = value.(float64)
	default:
		if value != nil {
			er = errors.New("Invalid ID type")
		}
	}
	if value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)
	}
	return val, err
}
