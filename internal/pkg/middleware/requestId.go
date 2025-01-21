package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/HeapSoil/auler/internal/pkg/utils"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求头中X-request-ID的新建或复用
		reqID := c.Request.Header.Get(utils.XRequestIDKey)
		if reqID == "" {
			// 新建32位
			reqID = uuid.New().String()
		}

		// 保存reqID在context中
		c.Set(utils.XRequestIDKey, reqID)

		// 保存reqID在HTTP返回头中
		c.Writer.Header().Set(utils.XRequestIDKey, reqID)
		c.Next()
	}
}
