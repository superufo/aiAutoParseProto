package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"game-backend/config"
)

// Claims JWT声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少认证令牌",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		// 解析JWT令牌
		claims, err := parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌无效: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 检查令牌是否过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（用于某些不需要强制登录的接口）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

// parseToken 解析JWT令牌
func parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	now := time.Now()
	expireTime := now.Add(time.Duration(config.AppConfig.JWT.ExpireTime) * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.AppConfig.JWT.Issuer,
			Subject:   string(rune(userID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	
	id, ok := userID.(uint)
	return id, ok
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	
	name, ok := username.(string)
	return name, ok
}

// GetClaims 从上下文中获取JWT声明
func GetClaims(c *gin.Context) (*Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	
	cl, ok := claims.(*Claims)
	return cl, ok
}
