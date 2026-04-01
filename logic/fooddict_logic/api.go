package fooddict_logic

import (
	"context"

	"github.com/MagicPig9898/familychefassistant_server/entity/fooddict_entity"
)

type FoodDictLogic interface {
	GetFoodDicts(ctx context.Context) ([]*fooddict_entity.TbFoodDict, error)

	GetAllFoodClass(ctx context.Context) ([]*fooddict_entity.TbFoodClass, error)

	AddFoodClass(ctx context.Context, class *fooddict_entity.TbFoodClass) error

	AddFoodDict(ctx context.Context, dict *fooddict_entity.TbFoodDict) error
}

func NewFoodDictLogic() FoodDictLogic {
	return newFoodDictLogicImpl()
}
