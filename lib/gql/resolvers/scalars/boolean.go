package scalars

import (
	"errors"
	"strconv"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type BoolScalar struct {
}

func NewBoolScalar() resolvers.Scalar {
	var scalar resolvers.Scalar
	scalar = &BoolScalar{}
	return scalar
}
func (o *BoolScalar) Assess(value interface{}) (val interface{}, err definitionError.GQLError) {
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
