package aenum

import (
	"reflect"
	"strconv"
)

// 新建一个基础接口,统一枚举类型的方法
type Enum interface {
	String() string
	In([]any) bool
}

// 基础结构体
type BaseEnum struct {
	value uint64
	name  string
}

func (e BaseEnum) String() string {
	return strconv.FormatUint(e.value, 10)
}

func (e BaseEnum) In(enums []any) bool {
	for _, enum := range enums {
		if e.value == reflect.ValueOf(enum).Uint() {
			return true
		}
	}
	return false
}
