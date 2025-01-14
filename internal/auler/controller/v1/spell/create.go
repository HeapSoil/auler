package spell

import (
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/utils"
	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// Create 用户创建一条咒语prompt
func (ctrl *SpellController) Create(c *gin.Context) {
	log.C(c).Infow("Create spell function called")

	var r v1.CreateSpellRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		errs.WriteResponse(c, errs.ErrBind, nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		errs.WriteResponse(c, errs.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	resp, err := ctrl.b.Spells().Create(c, c.GetString(utils.XUsernameKey), &r)
	if err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	errs.WriteResponse(c, nil, resp)

}
