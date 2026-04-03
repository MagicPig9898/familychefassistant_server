package fooddict_router

import (
	"net/http"

	"github.com/MagicPig9898/familychefassistant_server/entity/fooddict_entity"
	"github.com/MagicPig9898/familychefassistant_server/logic/fooddict_logic"
	"github.com/MagicPig9898/familychefassistant_server/router/router_utils"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	h := newHandler()
	registerGetFoodDicts(r, h)
	registerGetAllFoodClass(r, h)
	registerAddFoodClass(r, h)
	registerAddFoodDict(r, h)
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

func registerGetAllFoodClass(r *gin.RouterGroup, h *handler) {
	r.GET("/allfoodclass", func(c *gin.Context) {
		tbs, err := h.foodDictLogic.GetAllFoodClass(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, router_utils.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, router_utils.Success(tbs))
	})
}

func registerAddFoodClass(r *gin.RouterGroup, h *handler) {
	r.POST("/foodclass", func(c *gin.Context) {
		var class fooddict_entity.TbFoodClass
		if err := c.ShouldBindJSON(&class); err != nil {
			c.JSON(http.StatusBadRequest, router_utils.Fail(err.Error()))
			return
		}
		err := h.foodDictLogic.AddFoodClass(c.Request.Context(), &class)
		if err != nil {
			c.JSON(http.StatusInternalServerError, router_utils.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, router_utils.Success(nil))
	})
}

func registerAddFoodDict(r *gin.RouterGroup, h *handler) {
	r.POST("/fooddict", func(c *gin.Context) {
		var dict fooddict_entity.TbFoodDict
		if err := c.ShouldBindJSON(&dict); err != nil {
			c.JSON(http.StatusBadRequest, router_utils.Fail(err.Error()))
			return
		}
		err := h.foodDictLogic.AddFoodDict(c.Request.Context(), &dict)
		if err != nil {
			c.JSON(http.StatusInternalServerError, router_utils.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, router_utils.Success(nil))
	})
}
