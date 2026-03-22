package user_repo

import "context"

type UserRepo interface {
	GetUserByID(ctx context.Context, id int64) (string, error)
}

func NewUserRepo() UserRepo {
	return newUserRepoImpl()
}
