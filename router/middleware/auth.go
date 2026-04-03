package middleware

import (
	"net/http"
	"strings"

	jwt "github.com/MagicPig9898/familychefassistant_server/config/jwt_config"
	"github.com/MagicPig9898/familychefassistant_server/router/router_utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 返回一个 token 鉴权中间件
// whitelist 中的路径会直接放行，不检查 token
func AuthMiddleware(whitelist []string) gin.HandlerFunc {
	skip := make(map[string]struct{}, len(whitelist))
	for _, path := range whitelist {
		skip[path] = struct{}{}
	}

	return func(c *gin.Context) {
		// 白名单放行
		if _, ok := skip[c.FullPath()]; ok {
			c.Next()
			return
		}

		// 从 Header 中取 token：Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, router_utils.Fail("invalid_token"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, router_utils.Fail("invalid_token"))
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, router_utils.Fail("invalid_token"))
			c.Abort()
			return
		}

		// 将解析出的 open_id 存入上下文，后续 handler 可直接使用
		c.Set("open_id", claims.OpenID)
		c.Next()
	}
}
