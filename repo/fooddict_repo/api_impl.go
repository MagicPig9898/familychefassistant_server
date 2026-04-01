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
