package auler

import (
	"github.com/HeapSoil/auler/internal/auler/controller/v1/user"
	"github.com/HeapSoil/auler/internal/auler/store"
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/pkg/auth"
	"github.com/gin-gonic/gin"

	mw "github.com/HeapSoil/auler/internal/pkg/middleware"
)

func installRouters(g *gin.Engine) error {
	// 注册404 Handler
	g.NoRoute(func(c *gin.Context) {
		errs.WriteResponse(c, errs.ErrPageNotFound, nil)
	})

	// 注册/healthz Handler
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz  function called")
		errs.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return err
	}

	// user controller
	uc := user.New(store.S, authz)

	g.POST("/login", uc.Login)

	// 创建v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)                             // 创建用户
			userv1.PUT(":name/change-password", uc.ChangePassword) // 修改密码
			userv1.Use(mw.Authn(), mw.Authz(authz))                // 使用认证和鉴权中间件
			userv1.GET(":name", uc.Get)                            // 获取用户详情信息
		}
	}

	return nil

}
