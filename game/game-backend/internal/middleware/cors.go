package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		c.Header("Access-Control-Allow-Origin", "*")
		
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// 设置允许携带凭证
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// 设置预检请求的缓存时间
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
