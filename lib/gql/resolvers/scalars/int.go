package scalars

import (
	"errors"
	"strconv"

	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type IntScalar struct{
	
}

func NewIntScalar() resolvers.Scalar{
	var scalar resolvers.Scalar;
	scalar = &IntScalar{};
	return scalar
}
func(o *IntScalar) Assess(value interface{})( val interface{}, err error){
	var er error;

	switch value.(type){
	case string:
		val, er = strconv.ParseInt(value.(string),10, 64);
	case int ,int32 ,int64:
		val = value.(int64);
	default:
		if value != nil{
			err = errors.New("Invalid int type");
		}
	}
	if value != "" && er != nil{
		err = er;
	}
	return val, err
}
