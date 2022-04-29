package middlewares

import (
	"github.com/chenhg5/collection"

	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

var typeWhitelist = []string{"products/ssl", "tools/csr/decode"}

// AuthSsl ssl验证请求
func AuthSsl() gin.HandlerFunc {
	return func(c *gin.Context) {
		tag := c.DefaultQuery("type", "")
		if tag == "" {
			response.NormalVerificationError(c, "类型为空")
			c.Abort()
			return
		}
		if collection.Collect(typeWhitelist).Contains(tag) {
			c.Set("ssl_white_list", true)
		}
		c.Next()
	}
}
