package fooddict_repo

import (
	"context"

	"github.com/MagicPig9898/familychefassistant_server/entity/fooddict_entity"
)

type FoodDictRepo interface {
	GetFoodDicts(ctx context.Context) ([]*fooddict_entity.TbFoodDict, error)
}

func NewFoodDictRepo() FoodDictRepo {

	return newFoodDictRepoImpl()

}
