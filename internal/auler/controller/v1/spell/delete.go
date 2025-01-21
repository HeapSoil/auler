package spell

import (
	"github.com/gin-gonic/gin"

	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/utils"
)

func (ctrl *SpellController) Delete(c *gin.Context) {
	log.C(c).Infow("Delete spell function called")

	if err := ctrl.b.Spells().Delete(c, c.GetString(utils.XUsernameKey), c.Param("spellID")); err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	errs.WriteResponse(c, nil, nil)
}
