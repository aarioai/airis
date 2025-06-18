package dto

import "github.com/aarioai/airis/aa/atype"

type User struct {
	Uid      uint64 `json:"uid"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

type UserWithPaging struct {
	User
	atype.Paging
}
