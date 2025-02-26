package biz

import (
	"github.com/HeapSoil/auler/internal/auler/biz/spell"
	"github.com/HeapSoil/auler/internal/auler/biz/user"
	"github.com/HeapSoil/auler/internal/auler/store"
)

// 包级别的 Biz实例

type IBiz interface {
	Users() user.UserBiz
	Spells() spell.SpellBiz
}

type biz struct {
	ds store.IStore
}

var _ IBiz = (*biz)(nil)

func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}

func (b *biz) Spells() spell.SpellBiz {
	return spell.New(b.ds)
}
