package store

import (
	"sync"

	"gorm.io/gorm"
)

// 包级别的store实例

var (
	once sync.Once
	// 全局变量，方便其他包调用初始化好的S实例
	S *datastore
)

// Store层需要实现的方法interface
type IStore interface {
	Users() UserStore
}

// datastore是Istore的一个具体实现
type datastore struct {
	db *gorm.DB
}

// 用这个语句确保datastore有实现IStore的接口
var _ IStore = (*datastore)(nil)

// 全局New方法
func NewStore(db *gorm.DB) *datastore {
	// S只初始化一次
	once.Do(func() {
		S = &datastore{db}
	})
	return S
}

func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}
