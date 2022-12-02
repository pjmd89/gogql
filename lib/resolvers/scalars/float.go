package scalars

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/resolvers"
)

type Float float64

func NewFloatScalar() (r resolvers.Scalar) {
	var scalar *Float
	r = scalar
	return
}
func (o *Float) Set(value interface{}) (r interface{}, err definitionError.GQLError) {
	if value != nil {
		var rerr error
		s := fmt.Sprintf("%v", value)
		r, rerr = strconv.ParseFloat(s, 64)
		if rerr != nil {
			err = definitionError.NewFatal(rerr.Error(), nil)
		}
	}
	return
}
func (o *Float) Assess(resolved resolvers.ScalarResolved) (val interface{}, err definitionError.GQLError) {
	var er error

	switch resolved.Value.(type) {
	case string:
		val, er = strconv.ParseFloat(resolved.Value.(string), 64)
	case float32:
		val = float64(resolved.Value.(float32))
	case float64:
		val = float64(resolved.Value.(float64))
	case int:
		val = float64(resolved.Value.(int))
	case int32:
		val = float64(resolved.Value.(int32))
	case int64:
		val = float64(resolved.Value.(int64))
	default:
		if resolved.Value != nil {
			er = errors.New("Invalid float type")
		}
	}
	if resolved.Value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)
	}
	return val, err
}
