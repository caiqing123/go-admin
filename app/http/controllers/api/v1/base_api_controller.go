// Package v1 处理业务逻辑, api 控制器 v1
package v1

import (
	"reflect"
)

// BaseAPIController 基础控制器
type BaseAPIController struct {
}

//Transformation 转换处理
func Transformation(form, to interface{}) {
	fType := reflect.TypeOf(form)
	fValue := reflect.ValueOf(form)
	if fType.Kind() == reflect.Ptr {
		fType = fType.Elem()
		fValue = fValue.Elem()
	}
	tType := reflect.TypeOf(to)
	tValue := reflect.ValueOf(to)
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
		tValue = tValue.Elem()
	}
	for i := 0; i < fType.NumField(); i++ {
		for j := 0; j < tType.NumField(); j++ {
			if fType.Field(i).Name == tType.Field(j).Name &&
				fType.Field(i).Type.ConvertibleTo(tType.Field(j).Type) {
				tValue.Field(j).Set(fValue.Field(i))
			}
		}
	}
}
