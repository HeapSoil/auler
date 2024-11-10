package store

import (
	"context"

	"github.com/HeapSoil/auler/internal/pkg/model"
	"gorm.io/gorm"
)

// user模块在store层所实现的方法

type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
}

// users是UserStore接口的视线
type users struct {
	db *gorm.DB
}

var _ UserStore = (*users)(nil)

// 创建函数
func newUsers(db *gorm.DB) *users {
	return &users{db}
}

// 需要实现的方法
// Create 插入一条user记录

func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}
