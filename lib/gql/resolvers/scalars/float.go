package scalars

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type Float float64

func NewFloatScalar() (r resolvers.Scalar) {
	var scalar *Float
	r = scalar
	return
}
func (o *Float) Set(value interface{}) (r interface{}, err definitionError.GQLError) {
	s := fmt.Sprintf("%v", value)
	r, rerr := strconv.ParseFloat(s, 64)
	if rerr != nil {
		err = definitionError.NewFatal(rerr.Error(), nil)
	}
	return
}
func (o *Float) Assess(value interface{}) (val interface{}, err definitionError.GQLError) {
	var er error

	switch value.(type) {
	case string:
		val, er = strconv.ParseFloat(value.(string), 64)
	case float32:
		val = float64(value.(float32))
	case float64:
		val = float64(value.(float64))
	case int:
		val = float64(value.(int))
	case int32:
		val = float64(value.(int32))
	case int64:
		val = float64(value.(int64))
	default:
		if value != nil {
			er = errors.New("Invalid float type")
		}
	}
	if value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)
	}
	return val, err
}
