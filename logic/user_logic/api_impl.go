package userlogic

import (
	"context"

	userrepo "github.com/MagicPig9898/familychefassistant_server/repo/user_repo"
)
		
type userLogicImpl struct {
	repo userrepo.UserRepo
}

func newUserLogicImpl() *userLogicImpl {
	return &userLogicImpl{repo: userrepo.NewUserRepo()}
}

func (l *userLogicImpl) GetUserInfo(ctx context.Context, id int64) (string, error) {
	return l.repo.GetUserByID(ctx, id)
}
