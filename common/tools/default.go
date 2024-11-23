package tools

import (
	"errors"
	"reflect"
)

// Default 通过反射设置默认值
func Default(data any) error {
	typeOf := reflect.TypeOf(data)
	valueOf := reflect.ValueOf(data)

	if typeOf.Kind() != reflect.Pointer {
		return errors.New("param must be pointer")
	}

	ele := typeOf.Elem()
	valueEle := valueOf.Elem()

	for i := 0; i < ele.NumField(); i++ {
		field := ele.Field(i)
		value := valueEle.Field(i)

		kind := field.Type.Kind()
		if kind == reflect.Int {
			// 根据设置的tag进行值的设置
			value.Set(defaultInt())
		}
		if kind == reflect.Int32 {
			value.Set(defaultInt32())
		}
		if kind == reflect.Int64 {
			value.Set(defaultInt64())
		}
		if kind == reflect.Float32 {
			value.Set(defaultFloat32())
		}
		if kind == reflect.Float64 {
			value.Set(defaultFloat64())
		}
		if kind == reflect.String {
			value.Set(defaultString())
		}
	}
	return nil
}

func defaultInt() reflect.Value {
	var i = 0
	return reflect.ValueOf(i)
}

func defaultInt32() reflect.Value {
	var i int32 = 0
	return reflect.ValueOf(i)
}

func defaultInt64() reflect.Value {
	var i int64 = 0
	return reflect.ValueOf(i)
}

func defaultFloat32() reflect.Value {
	var i float32 = 0.0
	return reflect.ValueOf(i)
}

func defaultFloat64() reflect.Value {
	var i float64 = 0.0
	return reflect.ValueOf(i)
}

func defaultString() reflect.Value {
	var i = ""
	return reflect.ValueOf(i)
}
