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
			c.JSON(http.StatusBadRequest, router_utils.Fail("missing id"))
			return
		}
		id, err := router_utils.ParseInt64Query(idText)
		if err != nil {
			c.JSON(http.StatusBadRequest, router_utils.Fail("invalid id"))
			return
		}
		name, err := h.userLogic.GetUserInfo(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, router_utils.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusOK, router_utils.Success(gin.H{"id": id, "name": name})) // gin.H 就是 map[string]any 的类型别名，gin 会把它序列化成 JSON 返回给客户端
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
