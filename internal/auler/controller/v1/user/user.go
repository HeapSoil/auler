package user

import (
	"github.com/HeapSoil/auler/internal/auler/biz"
	"github.com/HeapSoil/auler/internal/auler/store"
	"github.com/HeapSoil/auler/pkg/auth"
)

type UserController struct {
	// 加入验证权限
	a *auth.Authz

	b biz.IBiz
}

func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{a: a, b: biz.NewBiz(ds)}
}
