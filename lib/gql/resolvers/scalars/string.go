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
	s := fmt.Sprintf("%v", value)
	r = s
	return
}
func (o *String) Assess(value interface{}) (val interface{}, err definitionError.GQLError) {
	var er error

	switch value.(type) {
	case string:
		val = value.(string)
	default:
		if value != nil {
			er = errors.New("Invalid string type")
		}
	}
	if value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)
	}
	return val, err
}
