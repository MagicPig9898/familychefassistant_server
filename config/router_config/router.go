package router_config

import (
	"github.com/MagicPig9898/familychefassistant_server/config/router_config/fooddict_router"
	"github.com/MagicPig9898/familychefassistant_server/config/router_config/user_router"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	Register(r)
	return r
}

func Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	user_router.Register(api.Group("/users"))
	fooddict_router.Register(api.Group("/fooddicts"))
}
