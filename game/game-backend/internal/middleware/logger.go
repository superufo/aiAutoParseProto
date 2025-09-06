package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 创建日志条目
		entry := logrus.WithFields(logrus.Fields{
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
			"status":      param.StatusCode,
			"latency":     param.Latency,
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"user_agent":  param.Request.UserAgent(),
			"error":       param.ErrorMessage,
		})

		// 根据状态码选择日志级别
		switch {
		case param.StatusCode >= 500:
			entry.Error("HTTP Request")
		case param.StatusCode >= 400:
			entry.Warn("HTTP Request")
		default:
			entry.Info("HTTP Request")
		}

		return ""
	})
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"path":  c.Request.URL.Path,
				"method": c.Request.Method,
			}).Error("Panic recovered")
		}

		c.JSON(500, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
	})
}
