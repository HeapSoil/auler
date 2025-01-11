package model

import (
	"time"

	"github.com/HeapSoil/auler/pkg/id"
	"gorm.io/gorm"
)

// SpellM 是数据库中 spell 记录 struct 格式的映射.
type SpellM struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Username  string    `gorm:"column:username;not null"`
	SpellID   string    `gorm:"column:spellID;not null"`
	Title     string    `gorm:"column:title;not null"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

// TableName 是 spell 映射的 MySQL 表名
func (s *SpellM) TableName() string {
	return "spell"
}

// BeforeCreate 使用 short-id 包在数据库记录生成前生成对应咒语的 id
func (s *SpellM) BeforeCreate(tx *gorm.DB) error {
	s.SpellID = "spell-" + id.GenShortID()

	return nil
}
