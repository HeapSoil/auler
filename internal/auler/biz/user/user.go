package user

import (
	"context"
	"regexp"

	"github.com/HeapSoil/auler/internal/auler/store"
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/model"
	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
	"github.com/jinzhu/copier"
)

// 定义接口，接口具体实现，并确保接口具体实现已经实现了定义接口
type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
}

type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

// New 创建一个实现了 UserBiz 接口的实例.
func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	// 采用copier简化代码
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		// 已经存在用户了
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errs.ErrUserAlreadyExists
		}

		return err
	}

	return nil

}
