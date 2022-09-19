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
	s := fmt.Sprintf("%v", value)
	r, rerr := strconv.ParseBool(s)
	if rerr != nil {
		err = definitionError.NewFatal(rerr.Error(), nil)
	}
	return
}
func (o *Bool) Assess(value interface{}) (val interface{}, err definitionError.GQLError) {
	var er error

	switch value.(type) {
	case string:
		val, er = strconv.ParseBool(value.(string))
	case bool:
		val = value.(bool)
	default:
		if value != nil {
			er = errors.New("Invalid bool type")

		}
	}
	if value != nil && er != nil {
		err = definitionError.NewWarning(er.Error(), nil)

	}
	return val, err
}
