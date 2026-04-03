package user_repo

import (
	"context"

	"github.com/MagicPig9898/easy_db/mysql"
	dbconfig "github.com/MagicPig9898/familychefassistant_server/config/db_config"
	"github.com/MagicPig9898/familychefassistant_server/entity/user_entity"
)

type userRepoImpl struct {
	mcli *mysql.Client
}

func newUserRepoImpl() *userRepoImpl {
	return &userRepoImpl{mcli: dbconfig.GetDb()}
}

func (r *userRepoImpl) GetUserByID(ctx context.Context, id string) (*user_entity.TbUser, error) {
	tb := &user_entity.TbUser{}
	sql := `select 
		id, 
		nick_name,
		avatar_url,
		city,
		province,
		country,
		gender,
		frist_login_time,
		last_login_time
	 from 
	 tb_user 
	 where id = ?`
	err := r.mcli.QueryOne(ctx, tb, sql, id)
	if err != nil {
		return nil, err
	}
	return tb, nil
}

func (r *userRepoImpl) InsertUser(ctx context.Context, user *user_entity.TbUser) error {
	sql := `insert into tb_user (
	    id,
		nick_name,
		avatar_url,
		city,
		province,
		country,
		gender,
		frist_login_time,
		last_login_time
	) values (
	 :id,
	 :nick_name, 
	 :avatar_url, 
	 :city, 
	 :province,
	 :country, 
	 :gender,
	 :frist_login_time,
	 :last_login_time)`
	err := r.mcli.InsertOneNamed(ctx, sql, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepoImpl) UpdateUser(ctx context.Context, user *user_entity.TbUser) error {
	sql := `update tb_user set 
		nick_name = :nick_name,
		avatar_url = :avatar_url,
		city = :city,
		province = :province,
		country = :country,
		gender = :gender,
		last_login_time = :last_login_time
	 where id = :id`
	err := r.mcli.UpdateOneNamed(ctx, sql, user)
	if err != nil {
		return err
	}
	return nil
}
