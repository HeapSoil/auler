package user


import (
	"github.com/gin-gonic/gin"

	"github.com/HeapSoil/auler/internal/pkg/log"
	"github.com/HeapSoil/auler/internal/pkg/core"
)

func (ctrl *UserController) Delete(c *gin.Context){
	log.C(c).Infow("Delete user function called")

	username := c.Param("name")

	if err := ctrl.b.Users().Delete(c, username); err != nil {
		core.WriteResponse(c, err, nil)
		return 
	}

	if _, err := ctrl.b.Users().Delete(c, username); err != nil {
		core.WriteResponse(c, err, nil)
		return 		
	}

	core.WriteResponse(c, nil, nil)

}