package scalars

import (
	"errors"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type IDScalar struct{
	
}
func NewIDScalar() resolvers.Scalar{
	var scalar resolvers.Scalar;
	scalar = &IDScalar{};
	return scalar
}
func(o *IDScalar) Assess(value interface{})( val interface{}, err definitionError.Error){
	var er error;

	switch value.(type){
	case string:
		val = value.(string);
	case int, int32, int64:
		val = value.(int64);
	case float32,float64:
		val = value.(float64);
	default:
		if value != nil{
			er = errors.New("Invalid ID type");
		}
	}
	if value != nil && er != nil{
		err=definitionError.NewWarning(er.Error(),nil);
	}
	return val, err
}
