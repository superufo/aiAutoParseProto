package model

import (
	"time"
	"gorm.io/gorm"
)

// Game 游戏模型
type Game struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	GameID      string         `json:"game_id" gorm:"uniqueIndex;size:50;not null"`
	Status      int            `json:"status" gorm:"default:0"` // 0:等待 1:进行中 2:已结束
	RoundID     string         `json:"round_id" gorm:"size:50"`
	Multiplier  float64        `json:"multiplier" gorm:"type:decimal(10,2);default:0"`
	PlayersCount int32         `json:"players_count" gorm:"default:0"`
	TotalBets   float64        `json:"total_bets" gorm:"type:decimal(15,2);default:0"`
	TotalPayout float64        `json:"total_payout" gorm:"type:decimal(15,2);default:0"`
	StartTime   *time.Time     `json:"start_time"`
	EndTime     *time.Time     `json:"end_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Bet 下注记录
type Bet struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	BetID        string         `json:"bet_id" gorm:"uniqueIndex;size:50;not null"`
	UserID       uint           `json:"user_id" gorm:"not null"`
	GameID       string         `json:"game_id" gorm:"size:50;not null"`
	RoundID      string         `json:"round_id" gorm:"size:50"`
	Amount       float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	AutoCashout  float64        `json:"auto_cashout" gorm:"type:decimal(10,2);default:0"`
	Multiplier   float64        `json:"multiplier" gorm:"type:decimal(10,2);default:0"`
	Payout       float64        `json:"payout" gorm:"type:decimal(15,2);default:0"`
	Status       int            `json:"status" gorm:"default:0"` // 0:进行中 1:已止盈 2:已崩盘
	CashoutTime  *time.Time     `json:"cashout_time"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// GameHistory 游戏历史记录
type GameHistory struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	RoundID       string    `json:"round_id" gorm:"uniqueIndex;size:50;not null"`
	GameID        string    `json:"game_id" gorm:"size:50;not null"`
	FinalMultiplier float64 `json:"final_multiplier" gorm:"type:decimal(10,2);not null"`
	PlayersCount  int32     `json:"players_count" gorm:"default:0"`
	TotalBets     float64   `json:"total_bets" gorm:"type:decimal(15,2);default:0"`
	TotalPayout   float64   `json:"total_payout" gorm:"type:decimal(15,2);default:0"`
	WinnersCount  int32     `json:"winners_count" gorm:"default:0"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	CreatedAt     time.Time `json:"created_at"`
}

// Leaderboard 排行榜
type Leaderboard struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	UserID            uint      `json:"user_id" gorm:"not null"`
	Username          string    `json:"username" gorm:"size:50;not null"`
	TotalWinnings     float64   `json:"total_winnings" gorm:"type:decimal(15,2);default:0"`
	BiggestMultiplier float64   `json:"biggest_multiplier" gorm:"type:decimal(10,2);default:0"`
	Rank              int32     `json:"rank" gorm:"default:0"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Game) TableName() string {
	return "games"
}

func (Bet) TableName() string {
	return "bets"
}

func (GameHistory) TableName() string {
	return "game_history"
}

func (Leaderboard) TableName() string {
	return "leaderboard"
}
