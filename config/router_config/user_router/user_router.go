package user_router

import (
	"net/http"

	"github.com/MagicPig9898/familychefassistant_server/config/router_config/router_utils"
	"github.com/MagicPig9898/familychefassistant_server/services/user_service"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	h := newHandler()
	registerHealthz(r)
	registerUserInfo(r, h)
}

type handler struct {
	userSvc user_service.UserService
}

func newHandler() *handler {
	return &handler{
		userSvc: user_service.NewUserService(),
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

		name, err := h.userSvc.GetUserInfo(c.Request.Context(), id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "name": name})
	})
}
