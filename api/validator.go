package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/indramhrt/simplebank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is supported or not
		return util.IsSupportedCurrenct(currency)
	}

	return false
}
