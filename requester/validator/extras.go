package validator

import (
	"net/mail"
	"reflect"
	"strconv"
	"time"
)

func isEmail(value reflect.Value) string {
	val, test := value.Interface().(string)

	if !test {
		return "must be a valid email"
	}
	_, err := mail.ParseAddress(val)
	if err != nil {
		return "must be a valid email"
	}
	return ""
}
func isDate(value reflect.Value) string {
	val, test := value.Interface().(string)
	if !test {
		return "must be a valid date"
	}
	layout := "02-01-2006"
	_, err := time.Parse(layout, val)
	if err != nil {
		return "must be a valid date"
	}
	return ""
}
func isNumericString(value reflect.Value) string {
	val, test := value.Interface().(string)
	if !test {
		return "a numeric string"
	}
	_, err := strconv.Atoi(val)
	if err != nil {
		return "a numeric string"
	}
	return ""
}

func isBooleanString(value reflect.Value) string {
	val, test := value.Interface().(string)
	if !test {
		return "must be a true or false string"
	}
	_, err := strconv.Atoi(val)
	if err != nil {
		return "must be a true or false string"
	}
	return ""
}
