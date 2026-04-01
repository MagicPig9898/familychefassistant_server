package fooddict_router

import (
	"net/http"

	"github.com/MagicPig9898/familychefassistant_server/config/router_config/router_utils"
	"github.com/MagicPig9898/familychefassistant_server/logic/fooddict_logic"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	h := newHandler()
	registerGetFoodDicts(r, h)
}

type handler struct {
	foodDictLogic fooddict_logic.FoodDictLogic
}

func newHandler() *handler {
	return &handler{
		foodDictLogic: fooddict_logic.NewFoodDictLogic(),
	}
}

func registerGetFoodDicts(r *gin.RouterGroup, h *handler) {
	r.GET("/fooddicts", func(c *gin.Context) {
		tbs, err := h.foodDictLogic.GetFoodDicts(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, router_utils.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, router_utils.Success(tbs))
	})

}
