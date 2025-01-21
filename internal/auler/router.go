package auler

import (
	"github.com/gin-gonic/gin"

	"github.com/HeapSoil/auler/internal/auler/controller/v1/spell"
	"github.com/HeapSoil/auler/internal/auler/controller/v1/user"
	"github.com/HeapSoil/auler/internal/auler/store"
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/pkg/auth"

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
	sc := spell.New(store.S)

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
			userv1.PUT(":name", uc.Update)                         // 更新用户
			userv1.GET("", uc.List)                                // 列出用户列表，只有 root 用户才能访问
			userv1.DELETE(":name", uc.Delete)                      // 删除用户
		}

		// spells 路由分组
		spellv1 := v1.Group("/spells", mw.Authn())
		{
			spellv1.POST("", sc.Create)             // 创建咒语
			spellv1.GET(":spellID", sc.Get)         // 获取咒语信息
			spellv1.PUT(":spellID", sc.Update)      // 更新咒语
			spellv1.GET("", sc.List)                // 获取咒语列表
			spellv1.DELETE(":spellID", sc.Delete)   // 删除咒语
			spellv1.DELETE("", sc.DeleteCollection) // 批量删除咒语
		}
	}

	return nil

}
