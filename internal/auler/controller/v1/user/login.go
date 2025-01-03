package user

import (
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
	"github.com/gin-gonic/gin"
)

func (ctrl *UserController) Login(c *gin.Context) {
	log.C(c).Infow("Login function called")

	var r v1.LoginRequest

	// BindJson出错
	if err := c.ShouldBindJSON(&r); err != nil {
		errs.WriteResponse(c, errs.ErrBind, nil)
		return
	}

	// 登陆业务出错
	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	errs.WriteResponse(c, nil, resp)
}
