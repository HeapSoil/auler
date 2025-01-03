package user

import (
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/asaskevich/govalidator"

	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
	"github.com/gin-gonic/gin"
)

const defaultMethods = "(GET)|(POST)|(PUT)|(DELETE)"

func (ctrl *UserController) Create(c *gin.Context) {
	log.C(c).Infow("Create user function called")

	var r v1.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		errs.WriteResponse(c, errs.ErrBind, nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		errs.WriteResponse(c, errs.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	if err := ctrl.b.Users().Create(c, &r); err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	// 加入authz认证
	if _, err := ctrl.a.AddNamedPolicy("p", r.Username, "/v1/users/"+r.Username, defaultMethods); err != nil {
		errs.WriteResponse(c, err, nil)

		return
	}

	errs.WriteResponse(c, nil, nil)

}
