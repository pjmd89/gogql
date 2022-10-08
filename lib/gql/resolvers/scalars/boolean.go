package scalars

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type Bool bool

func NewBoolScalar() (r resolvers.Scalar) {
	var scalar *Bool
	r = scalar
	return
}

func (o *Bool) Set(value interface{}) (r interface{}, err definitionError.GQLError) {
	if value != nil {
		var rerr error
		s := fmt.Sprintf("%v", value)
		r, rerr = strconv.ParseBool(s)
		if rerr != nil {
			err = definitionError.NewFatal(rerr.Error(), nil)
		}
	}
	return
}
func (o *Bool) Assess(resolved resolvers.ScalarResolved) (val interface{}, err definitionError.GQLError) {
	var er error

	switch resolved.Value.(type) {
	case string:
		val, er = strconv.ParseBool(resolved.Value.(string))
	case bool:
		val = resolved.Value.(bool)
	default:
		if resolved.Value != nil {
			er = errors.New("Invalid bool type")

		}
	}
	if resolved.Value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)

	}
	return val, err
}
