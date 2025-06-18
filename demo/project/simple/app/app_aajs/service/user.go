package service

import (
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"math/rand"
	"project/simple/app/app_aajs/module/bs/dto"
)

func (s *Service) Users(paging atype.Paging) ([]dto.UserWithPaging, *ae.Error) {
	users := make([]dto.UserWithPaging, 0, paging.PageSize)
	for i := 0; i < int(paging.PageSize); i++ {
		uid := uint64(paging.Offset) + uint64(i) + 1
		users = append(users, dto.UserWithPaging{
			User: dto.User{
				Uid:      uid,
				Username: fmt.Sprintf("User_%d", uid),
				Age:      rand.Intn(100),
			},
			Paging: paging,
		})
	}
	return users, nil
}

func (s *Service) User(uid uint64) (dto.User, *ae.Error) {
	user := dto.User{
		Uid:      uid,
		Username: fmt.Sprintf("User_%d", uid),
		Age:      rand.Intn(100),
	}
	return user, nil
}

func (s *Service) PostUser(username string, age int) (uint64, *ae.Error) {
	return rand.Uint64(), nil
}

func (s *Service) PutUser(uid uint64, username string, age int) *ae.Error {
	return nil
}

func (s *Service) PatchUser(uid uint64, age int) *ae.Error {
	return nil
}

func (s *Service) DeleteUser(uid uint64) *ae.Error {
	return nil
}
