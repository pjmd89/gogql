package scalars

import (
	"errors"
	"fmt"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type String string

func NewStringScalar() (r resolvers.Scalar) {
	var scalar *String
	r = scalar
	return
}
func (o *String) Set(value interface{}) (r interface{}, err definitionError.GQLError) {
	if value != nil {
		s := fmt.Sprintf("%v", value)
		r = s
	}
	return
}
func (o *String) Assess(resolved resolvers.ScalarResolved) (val interface{}, err definitionError.GQLError) {
	var er error

	switch resolved.Value.(type) {
	case string:
		val = resolved.Value.(string)
	default:
		if resolved.Value != nil {
			er = errors.New("Invalid string type")
		}
	}
	if resolved.Value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)
	}
	return val, err
}
