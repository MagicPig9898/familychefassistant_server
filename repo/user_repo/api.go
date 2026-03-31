package user_repo

import (
	"context"

	"github.com/MagicPig9898/familychefassistant_server/entity/user_entity"
)

type UserRepo interface {
	GetUserByID(ctx context.Context, id string) (*user_entity.TbUser, error)

	InsertUser(ctx context.Context, user *user_entity.TbUser) error

	UpdateUser(ctx context.Context, user *user_entity.TbUser) error
}

func NewUserRepo() UserRepo {
	return newUserRepoImpl()
}
