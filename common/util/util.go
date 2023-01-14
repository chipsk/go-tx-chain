package util

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"reflect"
)

func JsonString(data interface{}) string {
	str, _ := JsonStringWithError(data)
	return str
}

func JsonStringWithError(data interface{}) (string, error) {
	if IsNil(data) {
		return "", nil
	}

	typeKind := reflect.TypeOf(data)

	if typeKind == reflect.TypeOf("") {
		return data.(string), nil
	}

	// []byte, []uint8 特殊处理
	if typeKind == reflect.TypeOf([]uint8{}) {
		s := string(data.([]byte))
		if isJSON(s) {
			return s, nil
		}
	}

	jsonBytes, err := sonic.Marshal(data)
	if err != nil {
		describe := fmt.Sprintf("unable to marshal data as json string: +%v, +%v", err, data)
		return "", errors.New(describe)
	}

	res := string(jsonBytes)
	return res, nil
}

func IsNil(i interface{}) bool {
	k := reflect.ValueOf(i).Kind()
	if k == reflect.Struct {
		return reflect.ValueOf(&i).IsNil()
	}
	if k == reflect.Interface || k == reflect.Map || k == reflect.Slice || k == reflect.Ptr {
		return reflect.ValueOf(i).IsNil()
	}
	return true
}

func isJSON(s string) bool {
	var js interface{}
	return sonic.Unmarshal([]byte(s), &js) == nil
}
