package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"game-backend/internal/middleware"
	"game-backend/internal/model"
	"game-backend/internal/service"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token    string      `json:"token"`
	UserID   uint        `json:"user_id"`
	Username string      `json:"username"`
	Balance  float64     `json:"balance"`
	User     *model.User `json:"user"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证用户凭据
	user, err := h.authService.ValidateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 生成JWT Token
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成Token失败",
		})
		return
	}

	// 保存用户会话
	err = h.authService.SaveUserSession(user.ID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存会话失败",
		})
		return
	}

	response := LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Balance:  user.Balance,
		User:     user,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	exists, err := h.authService.CheckUsernameExists(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查用户名失败",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	exists, err = h.authService.CheckEmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查邮箱失败",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "邮箱已存在",
		})
		return
	}

	// 创建用户
	user, err := h.authService.CreateUser(req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "注册成功",
		"data": gin.H{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	// 删除用户会话
	err := h.authService.DeleteUserSession(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "登出失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登出成功",
	})
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    user,
	})
}

// UpdateProfile 更新用户信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	var req struct {
		Email  string `json:"email" binding:"omitempty,email"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 更新用户信息
	err := h.authService.UpdateUser(userID, req.Email, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// RefreshToken 刷新Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	username, exists := middleware.GetUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	// 生成新的Token
	newToken, err := middleware.GenerateToken(userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成Token失败",
		})
		return
	}

	// 更新用户会话
	err = h.authService.UpdateUserSession(userID, newToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新会话失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Token刷新成功",
		"data": gin.H{
			"token": newToken,
		},
	})
}
