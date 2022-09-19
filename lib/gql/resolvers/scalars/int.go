package scalars

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type Int int64

func NewIntScalar() (r resolvers.Scalar) {
	var scalar *Int
	r = scalar
	return
}
func (o *Int) Set(value interface{}) (r interface{}, err definitionError.GQLError) {
	s := fmt.Sprintf("%v", value)
	r, rerr := strconv.ParseInt(s, 10, 64)
	if rerr != nil {
		err = definitionError.NewFatal(rerr.Error(), nil)
	}
	return
}
func (o *Int) Assess(value interface{}) (val interface{}, err definitionError.GQLError) {
	var er error

	switch value.(type) {
	case string:
		val, er = strconv.ParseInt(value.(string), 10, 64)
	case int:
		val = value.(int)
	case int32:
		val = int(value.(int32))
	case int64:
		val = int(value.(int64))
	case float32:
		val = int(value.(float32))
	case float64:
		val = int(value.(float64))
	default:
		if value != nil {
			er = errors.New("Invalid int type")
		}
	}
	if value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)
	}
	return val, err
}
