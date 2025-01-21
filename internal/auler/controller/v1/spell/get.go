package spell

import (
	"github.com/gin-gonic/gin"

	"github.com/HeapSoil/auler/internal/pkg/errs"
	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/utils"
)

// 获取指定的咒语prompt
func (ctrl *SpellController) Get(c *gin.Context) {
	log.C(c).Infow("get spell function called")

	spell, err := ctrl.b.Spells().Get(c, c.GetString(utils.XUsernameKey), c.Param("spellID"))
	if err != nil {
		errs.WriteResponse(c, err, nil)
		return
	}

	errs.WriteResponse(c, nil, spell)
}
