package spell

import (
	"github.com/HeapSoil/auler/internal/auler/biz"
	"github.com/HeapSoil/auler/internal/auler/store"
)

// SpellController 是 spell 模块在 Controller 层的实现，用来处理咒语模块相关的请求.

type SpellController struct {
	b biz.IBiz
}

// New 新建 SpellController
func New(ds store.IStore) *SpellController {
	return &SpellController{
		b: biz.NewBiz(ds),
	}
}
