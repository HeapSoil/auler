package errs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	// Code 指定了业务错误码.
	Code string `json:"code"`

	// Message 包含了可以直接对外展示的错误信息.
	Message string `json:"message"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		// 解析error的错误，decode为三个状态变量
		// 根据错误类型，尝试从 err 中提取业务错误码和错误信息.
		httpCode, code, message := Decode(err)
		c.JSON(httpCode, ErrResponse{
			Code:    code,
			Message: message,
		})

		return

	}

	c.JSON(http.StatusOK, data)
}
