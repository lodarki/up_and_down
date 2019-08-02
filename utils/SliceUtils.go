package utils

import (
	"errors"
	"github.com/astaxie/beego"
	"reflect"
)

func SliceDel(arr interface{}, item interface{}) {
	val := reflect.ValueOf(arr)
	var value reflect.Value

	if val.Kind() == reflect.Ptr {
		value = reflect.Indirect(val)
	} else {
		value = val
	}

	if value.Kind() != reflect.Slice {
		beego.Error("invalid Param")
		return
	}

	var index = 0
	for index >= 0 {
		index = SliceIndexOf(arr, item)
		if index < 0 {
			return
		}

		oriLen := value.Len()
		if oriLen == 0 {
			return
		}
		var result []interface{}
		for i := 0; i < oriLen; i++ {
			if i != index {
				result = append(result, value.Index(i).Interface())
			}
		}

		value.SetLen(len(result))
		for i, v := range result {
			value.Index(i).Set(reflect.ValueOf(v))
		}
	}
}

func SliceIndexOf(arr interface{}, item interface{}) int {

	val := reflect.ValueOf(arr)
	var value reflect.Value

	if val.Kind() == reflect.Ptr {
		value = reflect.Indirect(val)
	} else {
		value = val
	}

	if value.Kind() != reflect.Slice {
		beego.Error("invalid Param")
		return -1
	}

	if value.Len() == 0 {
		return -1
	}

	for i := 0; i < value.Len(); i++ {
		if value.Index(i).Type() != reflect.TypeOf(item) {
			return -1
		}
		if value.Index(i).Interface() == reflect.ValueOf(item).Interface() {
			return i
		}
	}

	return -1
}

// 转置slice
func RevertSlice(arr interface{}) error {
	val := reflect.ValueOf(arr)
	if val.Kind() != reflect.Ptr {
		return errors.New("param must be a pointer, to make sure revaluing")
	}

	ind := reflect.Indirect(val)
	if ind.Kind() != reflect.Slice {
		return errors.New("invalid param")
	}

	var result = make([]interface{}, ind.Len())
	for i := 0; i < ind.Len(); i++ {
		result[ind.Len()-1-i] = ind.Index(i).Interface()
	}

	for j, v := range result {
		ind.Index(j).Set(reflect.ValueOf(v))
	}

	return nil
}

func SliceContains(arr interface{}, item interface{}) bool {
	return SliceIndexOf(arr, item) >= 0
}

func SliceSerialInt(ori []int, step int, scope int) (result []int) {

	if len(ori) == 0 {
		return
	}

	var count = 1
	for i, v := range ori {
		if i > 0 && v <= ori[i-1]+step {
			count += 1
			result = append(result, v)
		} else {
			count = 1
			result = []int{v}
		}

		if count >= scope {
			return
		}
	}

	result = []int{}
	return
}

func SliceSerialInt64(ori []int64, step int64, scope int64) (result []int64) {

	if len(ori) == 0 {
		return
	}

	var count int64 = 1
	for i, v := range ori {
		if i > 0 && v <= ori[i-1]+step {
			count += 1
			result = append(result, v)
		} else {
			count = 1
			result = []int64{v}
		}

		if count >= scope {
			return
		}
	}

	result = []int64{}
	return
}
