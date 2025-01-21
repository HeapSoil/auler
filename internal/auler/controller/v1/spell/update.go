package spell

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/utils"
	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
)

func (ctrl *SpellController) Update(c *gin.Context) {
	log.C(c).Infow("Update spell function called")

	var r v1.UpdateSpellRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		errs.WriteResponse(c, errs.ErrBind, nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		errs.WriteResponse(c, errs.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	if err := ctrl.b.Spells().Update(c, c.GetString(utils.XUsernameKey), c.Param("spellID"), &r); err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	errs.WriteResponse(c, nil, nil)
}
