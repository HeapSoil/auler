package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/utils"
)

// Auther 授权接口

type Auther interface {
	// 授权接口，其中：sub 操作主题， obj 操作对象，act 操作
	Authorize(sub, obj, act string) (bool, error)
}

// Authz 作为中间件从 gin.Context 解析用户名，用户请求路径，用户请求方法作为 casbin 授权模型进行鉴权
func Authz(a Auther) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetString(utils.XUsernameKey)
		obj := c.Request.URL.Path
		act := c.Request.Method

		log.Debugw("Build authorize context", "sub", sub, "obj", obj, "act", act)

		if allowed, _ := a.Authorize(sub, obj, act); !allowed {
			errs.WriteResponse(c, errs.ErrUnauthorized, nil)
			c.Abort()
			return
		}
	}
}
