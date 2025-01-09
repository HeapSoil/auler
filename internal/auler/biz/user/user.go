package user

import (
	"context"
	"errors"
	"regexp"

	"github.com/HeapSoil/auler/internal/auler/store"
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/model"
	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
	"github.com/HeapSoil/auler/pkg/auth"
	"github.com/HeapSoil/auler/pkg/token"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// 定义接口，接口具体实现，并确保接口具体实现已经实现了定义接口
type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error

	// 登陆业务
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	// 修改密码业务
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	// 获取用户信息业务（登陆后获取）
	Get(ctx context.Context, username string) (*v1.GetUserResponse, error)
	// 罗列用户业务
	List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error)
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

// 登陆业务
func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 先从请求中获取用户信息，再获取登这名录用户的所有信息
	user, err := b.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	matchErr := auth.Compare(user.Password, r.Password)
	if matchErr != nil {
		return nil, errs.ErrPasswordIncorrect
	}

	// 匹配成功，登陆成功，签发token返回
	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errs.ErrSignToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

// 修改密码业务逻辑
func (b *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	// 原来的密码和现有密码对不上
	if err := auth.Compare(userM.Password, r.OldPassword); err != nil {
		return errs.ErrPasswordIncorrect
	}

	// 加密密码存入db
	userM.Password, _ = auth.Encrypt(r.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil

}

// 获取用户信息业务（登陆后获取）
func (b *userBiz) Get(ctx context.Context, username string) (*v1.GetUserResponse, error) {
	user, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	var resp v1.GetUserResponse
	_ = copier.Copy(&resp, user)

	// mock data: 更新CreatedAt, UpdatedAt
	resp.CreatedAt = user.CreatedAt.Format("2024-10-02 17:14:55")
	resp.UpdatedAt = user.UpdatedAt.Format("2024-10-02 17:14:55")

	return &resp, nil
}

func (b *userBiz) List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(list))
	for _, item := range list {
		user := item
		users = append(users, &v1.UserInfo{
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	log.C(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &v1.ListUserResponse{TotalCount: count, Users: users}, nil

}
