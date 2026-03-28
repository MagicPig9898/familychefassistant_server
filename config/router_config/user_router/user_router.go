package user_router

import (
	"net/http"

	"github.com/MagicPig9898/familychefassistant_server/config/router_config/router_utils"
	"github.com/MagicPig9898/familychefassistant_server/logic/user_logic"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	h := newHandler()
	registerHealthz(r)
	registerUserInfo(r, h)
}

type handler struct {
	userLogic user_logic.UserLogic
}

func newHandler() *handler {
	return &handler{
		userLogic: user_logic.NewUserLogic(),
	}
}

func registerHealthz(r *gin.RouterGroup) {
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
}

func registerUserInfo(r *gin.RouterGroup, h *handler) {
	r.GET("/info", func(c *gin.Context) {
		idText := c.Query("id")
		if idText == "" {
			c.String(http.StatusBadRequest, "missing id")
			return
		}
		id, err := router_utils.ParseInt64Query(idText)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid id")
			return
		}
		name, err := h.userLogic.GetUserInfo(c.Request.Context(), id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "name": name})
	})
}
