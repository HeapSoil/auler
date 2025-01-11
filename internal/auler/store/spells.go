package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/HeapSoil/auler/internal/pkg/model"
)

// SpellStore 定义了 spell 模块在 store 层所实现的方法.
type SpellStore interface {
	Create(ctx context.Context, spell *model.SpellM) error
	Get(ctx context.Context, username, spellID string) (*model.SpellM, error)
	Update(ctx context.Context, spell *model.SpellM) error
	List(ctx context.Context, username string, offset, limit int) (int64, []*model.SpellM, error)
	Delete(ctx context.Context, username string, spellIDs []string) error
}

// SpellStore 接口的实现.
type spells struct {
	db *gorm.DB
}

// 确保 spells 实现了所有接口.
var _ SpellStore = (*spells)(nil)

func newSpells(db *gorm.DB) *spells {
	return &spells{db}
}

// Create 插入一条 spell 记录.
func (u *spells) Create(ctx context.Context, spell *model.SpellM) error {
	return u.db.Create(&spell).Error
}

// Get 根据 spellID 查询指定用户的 spell 数据库记录.
func (u *spells) Get(ctx context.Context, username, spellID string) (*model.SpellM, error) {
	var spell model.SpellM
	if err := u.db.Where("username = ? and spellID = ?", username, spellID).First(&spell).Error; err != nil {
		return nil, err
	}

	return &spell, nil
}

// Update 更新一条 spell 数据库记录.
func (u *spells) Update(ctx context.Context, spell *model.SpellM) error {
	return u.db.Save(spell).Error
}

// List 根据 offset 和 limit 返回指定用户的 spell 列表.
func (u *spells) List(ctx context.Context, username string, offset, limit int) (count int64, ret []*model.SpellM, err error) {
	err = u.db.Where("username = ?", username).Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

// Delete 根据 username, spellID 删除数据库 pospellst 记录.
func (u *spells) Delete(ctx context.Context, username string, spellIDs []string) error {
	err := u.db.Where("username = ? and spellID in (?)", username, spellIDs).Delete(&model.SpellM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
