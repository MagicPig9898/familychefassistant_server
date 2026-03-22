package userlogic

import "context"

type UserLogic interface {
	GetUserInfo(ctx context.Context, id int64) (string, error)
}

func NewUserLogic() UserLogic {
	return newUserLogicImpl()
}
