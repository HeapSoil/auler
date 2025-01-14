package spell

import (
	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (ctrl *SpellController) DeleteCollection(c *gin.Context) {
	log.C(c).Infow("delete spell collection function called")

	spellIDs := c.QueryArray("spellID")

	if err := ctrl.b.Spells().DeleteCollection(c, c.GetString(utils.XUsernameKey), spellIDs); err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	errs.WriteResponse(c, nil, nil)
}
