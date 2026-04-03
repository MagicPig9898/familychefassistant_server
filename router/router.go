package router

import (
	"github.com/MagicPig9898/familychefassistant_server/router/fooddict_router"
	"github.com/MagicPig9898/familychefassistant_server/router/middleware"
	"github.com/MagicPig9898/familychefassistant_server/router/user_router"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	Register(r)
	return r
}

func Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	// 注册 token 鉴权中间件，白名单内的路径不检查 token
	api.Use(middleware.AuthMiddleware([]string{
		"/api/v1/users/healthz",
		"/api/v1/users/wxlogin",
	}))

	user_router.Register(api.Group("/users"))
	fooddict_router.Register(api.Group("/fooddicts"))
}
