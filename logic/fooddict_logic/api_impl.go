package fooddict_logic

import (
	"context"

	"github.com/MagicPig9898/familychefassistant_server/entity/fooddict_entity"
	"github.com/MagicPig9898/familychefassistant_server/repo/fooddict_repo"
)

type foodDictLogicImpl struct {
	fdr fooddict_repo.FoodDictRepo
}

func newFoodDictLogicImpl() *foodDictLogicImpl {
	return &foodDictLogicImpl{fdr: fooddict_repo.NewFoodDictRepo()}
}

func (fd *foodDictLogicImpl) GetFoodDicts(ctx context.Context) ([]*fooddict_entity.TbFoodDict, error) {
	return fd.fdr.GetFoodDicts(ctx)
}
