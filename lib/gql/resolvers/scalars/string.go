package scalars

import (
	"errors"

	"github.com/pjmd89/gogql/lib/gql/definitionError"
	"github.com/pjmd89/gogql/lib/gql/resolvers"
)

type StringScalar struct{
	
}

func NewStringScalar() resolvers.Scalar{
	var scalar resolvers.Scalar;
	scalar = &StringScalar{};
	return scalar
}
func(o *StringScalar) Assess(value interface{})( val interface{}, err definitionError.Error){
	var er error;

	switch value.(type){
	case string:
		val = value.(string);
	default:
		if value != nil{
			er = errors.New("Invalid string type");
		}
	}
	if value != nil && er != nil{
		err=definitionError.NewWarning(er.Error(),nil);
	}
	return val, err
}
