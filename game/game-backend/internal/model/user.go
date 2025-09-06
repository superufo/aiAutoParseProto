package model

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"size:100"`
	Balance   float64        `json:"balance" gorm:"type:decimal(15,2);default:0"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Status    int            `json:"status" gorm:"default:1"` // 1:正常 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserStats 用户统计信息
type UserStats struct {
	ID                uint    `json:"id" gorm:"primaryKey"`
	UserID            uint    `json:"user_id" gorm:"not null"`
	TotalBets         int64   `json:"total_bets" gorm:"default:0"`
	TotalWinnings     float64 `json:"total_winnings" gorm:"type:decimal(15,2);default:0"`
	BiggestMultiplier float64 `json:"biggest_multiplier" gorm:"type:decimal(10,2);default:0"`
	GamesPlayed       int64   `json:"games_played" gorm:"default:0"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UserSession 用户会话
type UserSession struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;size:500;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (UserStats) TableName() string {
	return "user_stats"
}

func (UserSession) TableName() string {
	return "user_sessions"
}
