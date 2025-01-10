package user

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
	pb "github.com/HeapSoil/auler/pkg/proto/auler/v1"
)

// List 返回用户列表，只有 root 用户才能获取用户列表.
func (ctrl *UserController) List(c *gin.Context) {
	log.C(c).Infow("List user function called")

	var r v1.ListUserRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		errs.WriteResponse(c, errs.ErrBind, nil)
		return
	}

	resp, err := ctrl.b.Users().List(c, r.Offset, r.Limit)
	if err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	errs.WriteResponse(c, nil, resp)
}

// GRPC: ListUser 返回用户列表，只有 root 用户才能获取用户列表.
func (ctrl *UserController) ListUser(ctx context.Context, r *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	log.C(ctx).Infow("ListUser function called")

	resp, err := ctrl.b.Users().List(ctx, int(r.Offset), int(r.Limit))
	if err != nil {
		return nil, err
	}

	users := make([]*pb.UserInfo, 0, len(resp.Users))
	for _, u := range resp.Users {
		createdAt, _ := time.Parse("2006-01-02 15:04:05", u.CreatedAt)
		updatedAt, _ := time.Parse("2006-01-02 15:04:05", u.UpdatedAt)
		users = append(users, &pb.UserInfo{
			Username: u.Username,
			Nickname: u.Nickname,
			Email:    u.Email,
			Phone:    u.Phone,
			// SpellCount: u.SpellCount,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		})
	}

	ret := &pb.ListUserResponse{
		TotalCount: resp.TotalCount,
		Users:      users,
	}

	return ret, nil
}
