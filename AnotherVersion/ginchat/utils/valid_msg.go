package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
)

// GetValidMsg 返回结构体msg的参数
func GetValidMsg(err error, obj any) string {
	getObj := reflect.TypeOf(obj)
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		for _, e := range errs {
			// Elem() 方法用于获取指针指向的实际类型
			// 根据验证错误中的字段名（e.Field()）从结构体类型中查找对应的字段
			if f, exits := getObj.Elem().FieldByName(e.Field()); exits {
				msg := f.Tag.Get("msg")
				return msg
			}
		}
	}
	return err.Error()
}
