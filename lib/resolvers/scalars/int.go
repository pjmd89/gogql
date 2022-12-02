package scalars

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/resolvers"
)

type Int int64

func NewIntScalar() (r resolvers.Scalar) {
	var scalar *Int
	r = scalar
	return
}
func (o *Int) Set(value interface{}) (r interface{}, err definitionError.GQLError) {
	if value != nil {
		var rerr error
		s := fmt.Sprintf("%v", value)
		r, rerr = strconv.ParseInt(s, 10, 64)
		if rerr != nil {
			err = definitionError.NewFatal(rerr.Error(), nil)
		}
	}

	return
}
func (o *Int) Assess(resolved resolvers.ScalarResolved) (val interface{}, err definitionError.GQLError) {
	var er error

	switch resolved.Value.(type) {
	case string:
		val, er = strconv.ParseInt(resolved.Value.(string), 10, 64)
	case int:
		val = resolved.Value.(int)
	case int32:
		val = int(resolved.Value.(int32))
	case int64:
		val = int(resolved.Value.(int64))
	case float32:
		val = int(resolved.Value.(float32))
	case float64:
		val = int(resolved.Value.(float64))
	default:
		if resolved.Value != nil {
			er = errors.New("Invalid int type")
		}
	}
	if resolved.Value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)
	}
	return val, err
}
