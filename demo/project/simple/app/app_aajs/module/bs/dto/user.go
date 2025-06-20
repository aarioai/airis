package dto

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/atype/aenum"
)

type User struct {
	Uid      uint64    `json:"uid"`
	Username string    `json:"username"`
	Age      int       `json:"age"`
	Sex      aenum.Sex `json:"sex"`
}

type UserWithPaging struct {
	User
	atype.Paging
}
