package user_repo

import (
	"context"
	"fmt"

	"github.com/MagicPig9898/easy_db/mysql"
	dbconfig "github.com/MagicPig9898/familychefassistant_server/config/db_config"
)

type userRepoImpl struct {
	mcli *mysql.Client
}

func newUserRepoImpl() *userRepoImpl {
	dbcfg, err := dbconfig.GetDbConfig()
	if err != nil {
		fmt.Printf("Failed to get database clients: %v", err)
		return nil
	}
	return &userRepoImpl{mcli: dbcfg.MysqlCli}
}

func (r *userRepoImpl) GetUserByID(ctx context.Context, id int64) (string, error) {
	var name string
	err := r.mcli.QueryOne(ctx, &name, "SELECT name FROM users WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	return name, nil
}
