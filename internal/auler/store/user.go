package store

import (
	"context"

	"github.com/HeapSoil/auler/internal/pkg/model"
	"gorm.io/gorm"
)

// user模块在store层所实现的方法

type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	// 需要Get和Update获取用户信息（如密码对比），并修改用户信息（如密码）
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
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

// Get 根据用户名查询指定 user 的数据库记录.
func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新一条 user 数据库记录
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}
