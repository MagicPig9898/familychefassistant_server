package user_repo

import (
	"context"

	"github.com/MagicPig9898/easy_db/mysql"
	dbconfig "github.com/MagicPig9898/familychefassistant_server/config/db_config"
)

type userRepoImpl struct {
	mcli *mysql.Client
}

func newUserRepoImpl() *userRepoImpl {
	return &userRepoImpl{mcli: dbconfig.GetDb()}
}

func (r *userRepoImpl) GetUserByID(ctx context.Context, id int64) (string, error) {
	var name string
	err := r.mcli.QueryOne(ctx, &name, "SELECT name FROM users WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	return name, nil
}
