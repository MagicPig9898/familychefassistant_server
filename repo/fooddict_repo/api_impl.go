package fooddict_repo

import (
	"context"

	"github.com/MagicPig9898/easy_db/mysql"
	dbconfig "github.com/MagicPig9898/familychefassistant_server/config/db_config"
	log "github.com/MagicPig9898/familychefassistant_server/config/log_config"
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
		log.Errorf("GetFoodDicts failed: %v", err)
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
		log.Errorf("GetAllFoodClass failed: %v", err)
		return nil, err
	}
	return tbs, nil
}

func (fd *foodDictRepoImpl) AddFoodClass(ctx context.Context, class *fooddict_entity.TbFoodClass) error {
	sql := `insert into tb_food_class (id, name) values (:id, :name)`
	err := fd.mcli.InsertOneNamed(ctx, sql, class)
	if err != nil {
		log.Errorf("AddFoodClass failed: %v", err)
		return err
	}
	return nil
}

func (fd *foodDictRepoImpl) AddFoodDict(ctx context.Context, dict *fooddict_entity.TbFoodDict) error {
	sql := `insert into tb_food_dict (id, class_id, name, image, description, good, bad) values (:id, :class_id, :name, :image, :description)`
	err := fd.mcli.InsertOneNamed(ctx, sql, dict)
	if err != nil {
		log.Errorf("AddFoodDict failed: %v", err)
		return err
	}
	return nil
}
