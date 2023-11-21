package validation

import (
	"fmt"
	"go-rest-api/adapter/validator"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ozzo struct {
	err error
	msg []string
}

func NewOzzo() (validator.Validator, error) {
	return &ozzo{}, nil
}

func (o *ozzo) Validate(i interface{}) error {
	if len(o.msg) > 0 {
		o.msg = nil
	}

	o.err = validation.Validate(i)

	return o.err
}

func (o ozzo) Messages() []string {
	if ve, ok := o.err.(validation.Errors); ok {
		var errorList []string
		for field, e := range ve {
			errorList = append(errorList, fmt.Sprintf("%v: %v", field, e))
		}
		return errorList
	}
	return []string{o.err.Error()}
	// bson, err := json.Marshal(o.err)
	// if err != nil {
	// 	if ve, ok := err.(validation.Errors); ok {
	// 		errMap := make(map[string]string)
	// 		for field, e := range ve {
	// 			errMap[field] = e.Error()
	// 		}
	// 		jsonErr, err := json.Marshal(errMap)
	// 		if err != nil {
	// 			return o.err.Error()
	// 		}
	// 		return string(jsonErr)
	// 	}
	// }
	// return string(bson)
}
