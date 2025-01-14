package spell

import (
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/utils"
	v1 "github.com/HeapSoil/auler/pkg/api/auler/v1"
	"github.com/gin-gonic/gin"
)

func (ctrl *SpellController) List(c *gin.Context) {
	log.C(c).Infow("List spells function called")

	var r v1.ListSpellRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		errs.WriteResponse(c, errs.ErrBind, nil)
		return
	}

	resp, err := ctrl.b.Spells().List(c, c.GetString(utils.XUsernameKey), r.Offset, r.Limit)
	if err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	errs.WriteResponse(c, nil, resp)
}
