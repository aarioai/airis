package afmt

import (
	"fmt"
	"reflect"
)

// ErrmsgSideEffect 对于不应该修改入参指针的（一般会备注 @note slice 入参安全），函数修改了后返回的错误提示。
// 一般用于函数或单元测试错误提示
func ErrmsgSideEffect(variable any) string {
	return fmt.Sprintf("improper use of %s occurs unintended side effects", reflect.TypeOf(variable))
}
