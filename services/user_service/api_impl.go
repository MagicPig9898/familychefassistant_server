package user_service

import (
	"context"

	user_logic "github.com/MagicPig9898/familychefassistant_server/logic/user_logic"
)

type userServiceImpl struct {
	repo user_logic.UserLogic
}

func newUserServiceImpl() *userServiceImpl {
	return &userServiceImpl{repo: user_logic.NewUserLogic()}
}

func (s *userServiceImpl) GetUserInfo(ctx context.Context, id int64) (string, error) {
	return s.repo.GetUserInfo(ctx, id)
}
