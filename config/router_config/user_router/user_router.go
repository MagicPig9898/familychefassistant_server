package user_router

import (
	"net/http"

	"github.com/MagicPig9898/familychefassistant_server/config/router_config/router_utils"
	"github.com/MagicPig9898/familychefassistant_server/entity/user_entity"
	"github.com/MagicPig9898/familychefassistant_server/logic/user_logic"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	h := newHandler()
	registerHealthz(r)
	registerUserInfo(r, h)
	registerWXLogin(r, h)
	registerValidToken(r, h)
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
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, router_utils.Fail("missing id"))
			return
		}
		tb, err := h.userLogic.GetUserInfo(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, router_utils.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, router_utils.Success(tb))
	})
}

func registerWXLogin(r *gin.RouterGroup, h *handler) {
	r.POST("/wxlogin", func(c *gin.Context) {
		var req user_entity.UserLoginDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, router_utils.Fail("参数错误"))
			return
		}
		resp, err := h.userLogic.WXLogin(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, router_utils.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, router_utils.Success(resp))
	})
}

func registerValidToken(r *gin.RouterGroup, h *handler) {
	r.POST("/validtoken", func(c *gin.Context) {
		var req user_entity.ValidTokenDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, router_utils.Fail("参数错误"))
			return
		}
		newToken, err := h.userLogic.ValidToken(c.Request.Context(), req.Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, router_utils.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, router_utils.Success(newToken))
	})
}
