package fooddict_logic

import (
	"context"

	"github.com/MagicPig9898/familychefassistant_server/entity/fooddict_entity"
)

type FoodDictLogic interface {
	GetFoodDicts(ctx context.Context) ([]*fooddict_entity.TbFoodDict, error)
}

func NewFoodDictLogic() FoodDictLogic {

	return newFoodDictLogicImpl()

}
