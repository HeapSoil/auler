package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/utils"
	"github.com/HeapSoil/auler/pkg/token"
)

// 认证方法Authn 中间件
func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析 JWT Token
		username, err := token.ParseRequest(c)
		if err != nil {
			errs.WriteResponse(c, errs.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Set(utils.XUsernameKey, username)
		// 下一个中间件
		c.Next()
	}
}
