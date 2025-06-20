package service

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/atype/aenum"
	"github.com/aarioai/airis/pkg/basic"
	"math/rand"
	"project/simple/app/app_aajs/module/bs/dto"
)

func (s *Service) Users(ctx context.Context, paging atype.Paging) ([]dto.UserWithPaging, *ae.Error) {
	users := make([]dto.UserWithPaging, 0, paging.PageSize)
	for i := 0; i < int(paging.PageSize); i++ {
		uid := uint64(paging.Offset) + uint64(i) + 1
		users = append(users, dto.UserWithPaging{
			User: dto.User{
				Uid:      uid,
				Username: fmt.Sprintf("User_%d", uid),
				Age:      rand.Intn(100),
				Sex:      basic.Ter(uid%2 == 0, aenum.Female, aenum.Male),
			},
			Paging: paging,
		})
	}
	return users, nil
}

func (s *Service) QueryUsersBySex(ctx context.Context, sex aenum.Sex, paging atype.Paging) ([]dto.UserWithPaging, *ae.Error) {
	users := make([]dto.UserWithPaging, 0, paging.PageSize)
	for i := 0; i < int(paging.PageSize); i++ {
		uid := uint64(paging.Offset*2) + uint64(i*2) + uint64(sex)
		users = append(users, dto.UserWithPaging{
			User: dto.User{
				Uid:      uid,
				Username: fmt.Sprintf("User_%d", uid),
				Age:      rand.Intn(100),
				Sex:      basic.Ter(uid%2 == 0, aenum.Female, aenum.Male),
			},
			Paging: paging,
		})
	}
	return users, nil
}

func (s *Service) QueryUsers(ctx context.Context, uids []uint64) ([]dto.User, *ae.Error) {
	users := make([]dto.User, 0, len(uids))
	for _, uid := range uids {
		users = append(users, dto.User{
			Uid:      uid,
			Username: fmt.Sprintf("User_%d", uid),
			Age:      rand.Intn(100),
			Sex:      basic.Ter(uid%2 == 0, aenum.Female, aenum.Male),
		})
	}
	return users, nil
}

func (s *Service) User(ctx context.Context, uid uint64) (dto.User, *ae.Error) {
	user := dto.User{
		Uid:      uid,
		Username: fmt.Sprintf("User_%d", uid),
		Age:      rand.Intn(100),
		Sex:      basic.Ter(uid%2 == 0, aenum.Female, aenum.Male),
	}
	return user, nil
}

func (s *Service) PostUser(ctx context.Context, username string, age int) (uint64, *ae.Error) {
	return rand.Uint64(), nil
}

func (s *Service) PutUser(ctx context.Context, uid uint64, username string, age int) *ae.Error {
	return nil
}

func (s *Service) PatchUser(ctx context.Context, uid uint64, age int) *ae.Error {
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, uid uint64) *ae.Error {
	return nil
}
