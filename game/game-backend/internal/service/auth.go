package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"game-backend/internal/model"
	"game-backend/config"
)

// AuthService 认证服务
type AuthService struct {
	db *gorm.DB
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

// ValidateUser 验证用户凭据
func (s *AuthService) ValidateUser(username, password string) (*model.User, error) {
	var user model.User
	
	// 查找用户
	if err := s.db.Where("username = ? AND status = 1", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

// CheckUsernameExists 检查用户名是否存在
func (s *AuthService) CheckUsernameExists(username string) (bool, error) {
	var count int64
	err := s.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// CheckEmailExists 检查邮箱是否存在
func (s *AuthService) CheckEmailExists(email string) (bool, error) {
	var count int64
	err := s.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// CreateUser 创建用户
func (s *AuthService) CreateUser(username, password, email string) (*model.User, error) {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Balance:  100.0, // 新用户初始余额
		Status:   1,     // 正常状态
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// 创建用户统计记录
	userStats := &model.UserStats{
		UserID:            user.ID,
		TotalBets:         0,
		TotalWinnings:     0,
		BiggestMultiplier: 0,
		GamesPlayed:       0,
	}

	if err := s.db.Create(userStats).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// SaveUserSession 保存用户会话
func (s *AuthService) SaveUserSession(userID uint, token string) error {
	// 删除旧会话
	s.db.Where("user_id = ?", userID).Delete(&model.UserSession{})

	// 创建新会话
	session := &model.UserSession{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireTime) * time.Hour),
	}

	return s.db.Create(session).Error
}

// DeleteUserSession 删除用户会话
func (s *AuthService) DeleteUserSession(userID uint) error {
	return s.db.Where("user_id = ?", userID).Delete(&model.UserSession{}).Error
}

// UpdateUserSession 更新用户会话
func (s *AuthService) UpdateUserSession(userID uint, token string) error {
	session := &model.UserSession{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireTime) * time.Hour),
	}

	return s.db.Where("user_id = ?", userID).Updates(session).Error
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := s.db.Where("id = ? AND status = 1", userID).First(&user).Error
	return &user, err
}

// UpdateUser 更新用户信息
func (s *AuthService) UpdateUser(userID uint, email, avatar string) error {
	updates := make(map[string]interface{})
	
	if email != "" {
		updates["email"] = email
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}

	if len(updates) == 0 {
		return errors.New("没有需要更新的字段")
	}

	return s.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

// ValidateToken 验证Token
func (s *AuthService) ValidateToken(token string) (*model.User, error) {
	var session model.UserSession
	
	// 查找会话
	if err := s.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("会话不存在或已过期")
		}
		return nil, err
	}

	// 获取用户信息
	return s.GetUserByID(session.UserID)
}

// CleanExpiredSessions 清理过期会话
func (s *AuthService) CleanExpiredSessions() error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&model.UserSession{}).Error
}
