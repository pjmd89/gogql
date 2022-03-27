package scalars

import (
	"errors"
	"strconv"

	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type FloatScalar struct{
	
}

func NewFloatScalar() resolvers.Scalar{
	var scalar resolvers.Scalar;
	scalar = &FloatScalar{};
	return scalar
}
func(o *FloatScalar) Assess(value interface{})( val interface{}, err error){
	var er error;

	switch value.(type){
	case string:
		val, er = strconv.ParseFloat(value.(string), 64);
	case float32, float64 ,int ,int32 ,int64:
		val = value.(float64);
	default:
		if value != nil{
			err = errors.New("Invalid float type");
		}
	}
	if value != "" && er != nil{
		err = er;
	}
	return val, err
}
