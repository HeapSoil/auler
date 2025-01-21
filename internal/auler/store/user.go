package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/HeapSoil/auler/internal/pkg/model"
)

// user模块在store层所实现的方法

type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	// 需要Get和Update获取用户信息（如密码对比），并修改用户信息（如密码）
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	// 需要罗列用户方法，返回个数和每个用户
	List(ctx context.Context, offset, limit int) (int64, []*model.UserM, error)

	// 管理员删除用户
	Delete(ctx context.Context, username string) error
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

func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*model.UserM, err error) {
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

// Delete 根据 username 删除数据库对应的用户记录
func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&model.UserM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
