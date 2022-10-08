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
	if value != nil {
		s := fmt.Sprintf("%v", value)
		r = s
	}

	return
}
func (o *ID) Assess(resolved resolvers.ScalarResolved) (val interface{}, err definitionError.GQLError) {
	var er error

	switch resolved.Value.(type) {
	case string:
		val = resolved.Value.(string)
	case int, int32, int64:
		val = resolved.Value.(int64)
	case float32, float64:
		val = resolved.Value.(float64)
	default:
		if resolved.Value != nil {
			er = errors.New("Invalid ID type")
		}
	}
	if resolved.Value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)
	}
	return val, err
}
