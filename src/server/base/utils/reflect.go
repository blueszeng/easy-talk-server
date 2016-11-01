package utils

import (
	"reflect"
)

func IsNil(i interface{}) bool {
	defer Recover()

	if i == nil {
		return true
	}

	return reflect.ValueOf(i).IsNil()
}
