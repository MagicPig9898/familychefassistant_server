package user_logic

import (
	"context"

	"github.com/MagicPig9898/familychefassistant_server/entity/user_entity"
)

type UserLogic interface {
	GetUserInfo(ctx context.Context, id string) (*user_entity.TbUser, error)

	WXLogin(ctx context.Context, userLoginDto *user_entity.UserLoginDto) (*user_entity.UserLoginDto, error)

	ValidToken(ctx context.Context, token string) (string, error)
}

func NewUserLogic() UserLogic {
	return newUserLogicImpl()
}
