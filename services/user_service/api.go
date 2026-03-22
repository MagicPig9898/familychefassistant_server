package user_service

import "context"

type UserService interface {
	GetUserInfo(ctx context.Context, id int64) (string, error)
}

func NewUserService() UserService {
	return newUserServiceImpl()
}
