package fooddict_repo

import (
	"context"

	"github.com/MagicPig9898/easy_db/mysql"
	dbconfig "github.com/MagicPig9898/familychefassistant_server/config/db_config"
	"github.com/MagicPig9898/familychefassistant_server/entity/fooddict_entity"
)

type foodDictRepoImpl struct {
	mcli *mysql.Client
}

func newFoodDictRepoImpl() *foodDictRepoImpl {
	return &foodDictRepoImpl{mcli: dbconfig.GetDb()}
}

func (fd *foodDictRepoImpl) GetFoodDicts(ctx context.Context) ([]*fooddict_entity.TbFoodDict, error) {
	tbs := []*fooddict_entity.TbFoodDict{}
	sql := `select 
		*
	 from 
	    tb_food_dict`
	err := fd.mcli.QueryMany(ctx, &tbs, sql)

	if err != nil {
		return nil, err
	}
	return tbs, nil
}

func (fd *foodDictRepoImpl) GetAllFoodClass(ctx context.Context) ([]*fooddict_entity.TbFoodClass, error) {
	tbs := []*fooddict_entity.TbFoodClass{}
	sql := `select 
		*
	 from 
	    tb_food_class`
	err := fd.mcli.QueryMany(ctx, &tbs, sql)
	if err != nil {
		return nil, err
	}
	return tbs, nil
}

func (fd *foodDictRepoImpl) AddFoodClass(ctx context.Context, class *fooddict_entity.TbFoodClass) error {
	sql := `insert into tb_food_class (id, name) values (:id, :name)`
	err := fd.mcli.InsertOneNamed(ctx, sql, class)
	if err != nil {
		return err
	}
	return nil
}

func (fd *foodDictRepoImpl) AddFoodDict(ctx context.Context, dict *fooddict_entity.TbFoodDict) error {
	sql := `insert into tb_food_dict (id, class_id, name, image, description, good, bad) values (:id, :class_id, :name, :image, :description, :good, :bad)`
	err := fd.mcli.InsertOneNamed(ctx, sql, dict)
	if err != nil {
		return err
	}
	return nil
}
