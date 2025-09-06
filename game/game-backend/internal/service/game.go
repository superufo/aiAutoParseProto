package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"game-backend/internal/model"
	"game-backend/config"
)

// GameService 游戏服务
type GameService struct {
	db *gorm.DB
}

// NewGameService 创建游戏服务
func NewGameService(db *gorm.DB) *GameService {
	return &GameService{
		db: db,
	}
}

// CreateBet 创建下注记录
func (s *GameService) CreateBet(userID uint, amount, autoCashout float64) (*model.Bet, error) {
	betID := fmt.Sprintf("bet_%d_%d", userID, time.Now().Unix())
	
	bet := &model.Bet{
		BetID:       betID,
		UserID:      userID,
		GameID:      "crash_001",
		Amount:      amount,
		AutoCashout: autoCashout,
		Status:      0, // 进行中
	}

	if err := s.db.Create(bet).Error; err != nil {
		return nil, err
	}

	return bet, nil
}

// GetBetByID 根据ID获取下注记录
func (s *GameService) GetBetByID(betID string) (*model.Bet, error) {
	var bet model.Bet
	err := s.db.Where("bet_id = ?", betID).First(&bet).Error
	return &bet, err
}

// UpdateBetCashout 更新下注止盈信息
func (s *GameService) UpdateBetCashout(betID string, multiplier, payout float64) error {
	now := time.Now()
	return s.db.Model(&model.Bet{}).Where("bet_id = ?", betID).Updates(map[string]interface{}{
		"multiplier":   multiplier,
		"payout":       payout,
		"status":       1, // 已止盈
		"cashout_time": &now,
	}).Error
}

// UpdateBetCrash 更新下注崩盘信息
func (s *GameService) UpdateBetCrash(betID string) error {
	return s.db.Model(&model.Bet{}).Where("bet_id = ?", betID).Updates(map[string]interface{}{
		"status": 2, // 已崩盘
	}).Error
}

// GetUserBetHistory 获取用户下注历史
func (s *GameService) GetUserBetHistory(userID uint, page, pageSize int) ([]model.Bet, int64, error) {
	var bets []model.Bet
	var total int64

	// 获取总数
	if err := s.db.Model(&model.Bet{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bets).Error

	return bets, total, err
}

// GetGameHistory 获取游戏历史
func (s *GameService) GetGameHistory(page, pageSize int) ([]model.GameHistory, int64, error) {
	var games []model.GameHistory
	var total int64

	// 获取总数
	if err := s.db.Model(&model.GameHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := s.db.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&games).Error

	return games, total, err
}

// GetLeaderboard 获取排行榜
func (s *GameService) GetLeaderboard() ([]model.Leaderboard, error) {
	var leaderboard []model.Leaderboard
	
	err := s.db.Order("total_winnings DESC").
		Limit(100).
		Find(&leaderboard).Error

	return leaderboard, err
}

// GetUserStats 获取用户统计
func (s *GameService) GetUserStats(userID uint) (*model.UserStats, error) {
	var stats model.UserStats
	err := s.db.Where("user_id = ?", userID).First(&stats).Error
	return &stats, err
}

// UpdateUserStats 更新用户统计
func (s *GameService) UpdateUserStats(userID uint, betAmount, payout, multiplier float64) error {
	stats := &model.UserStats{}
	
	// 查找或创建统计记录
	if err := s.db.Where("user_id = ?", userID).First(stats).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新统计记录
			stats = &model.UserStats{
				UserID:            userID,
				TotalBets:         1,
				TotalWinnings:     payout,
				BiggestMultiplier: multiplier,
				GamesPlayed:       1,
			}
			return s.db.Create(stats).Error
		}
		return err
	}

	// 更新统计信息
	updates := map[string]interface{}{
		"total_bets":     stats.TotalBets + 1,
		"total_winnings": stats.TotalWinnings + payout,
		"games_played":   stats.GamesPlayed + 1,
	}

	// 更新最大倍数
	if multiplier > stats.BiggestMultiplier {
		updates["biggest_multiplier"] = multiplier
	}

	return s.db.Model(stats).Updates(updates).Error
}

// DeductUserBalance 扣除用户余额
func (s *GameService) DeductUserBalance(userID uint, amount float64) error {
	return s.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("balance", gorm.Expr("balance - ?", amount)).Error
}

// AddUserBalance 增加用户余额
func (s *GameService) AddUserBalance(userID uint, amount float64) error {
	return s.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

// GetUserByID 根据ID获取用户
func (s *GameService) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := s.db.Where("id = ? AND status = 1", userID).First(&user).Error
	return &user, err
}

// CreateGameHistory 创建游戏历史记录
func (s *GameService) CreateGameHistory(roundID, gameID string, finalMultiplier float64, playersCount int32, totalBets, totalPayout float64, winnersCount int32) error {
	history := &model.GameHistory{
		RoundID:         roundID,
		GameID:          gameID,
		FinalMultiplier: finalMultiplier,
		PlayersCount:    playersCount,
		TotalBets:       totalBets,
		TotalPayout:     totalPayout,
		WinnersCount:    winnersCount,
		StartTime:       time.Now().Add(-30 * time.Second), // 假设游戏持续30秒
		EndTime:         time.Now(),
	}

	return s.db.Create(history).Error
}

// UpdateLeaderboard 更新排行榜
func (s *GameService) UpdateLeaderboard() error {
	// 删除旧排行榜数据
	if err := s.db.Exec("DELETE FROM leaderboard").Error; err != nil {
		return err
	}

	// 重新计算排行榜
	sql := `
		INSERT INTO leaderboard (user_id, username, total_winnings, biggest_multiplier, rank, updated_at)
		SELECT 
			u.id as user_id,
			u.username,
			COALESCE(us.total_winnings, 0) as total_winnings,
			COALESCE(us.biggest_multiplier, 0) as biggest_multiplier,
			ROW_NUMBER() OVER (ORDER BY COALESCE(us.total_winnings, 0) DESC) as rank,
			NOW() as updated_at
		FROM users u
		LEFT JOIN user_stats us ON u.id = us.user_id
		WHERE u.status = 1
		ORDER BY COALESCE(us.total_winnings, 0) DESC
		LIMIT 100
	`

	return s.db.Exec(sql).Error
}

// GetMinBetAmount 获取最小下注金额
func (s *GameService) GetMinBetAmount() float64 {
	return config.AppConfig.Game.MinBetAmount
}

// GetMaxBetAmount 获取最大下注金额
func (s *GameService) GetMaxBetAmount() float64 {
	return config.AppConfig.Game.MaxBetAmount
}

// GetMinMultiplier 获取最小倍数
func (s *GameService) GetMinMultiplier() float64 {
	return config.AppConfig.Game.MinMultiplier
}

// GetMaxMultiplier 获取最大倍数
func (s *GameService) GetMaxMultiplier() float64 {
	return config.AppConfig.Game.MaxMultiplier
}

// ProcessAutoCashouts 处理自动止盈
func (s *GameService) ProcessAutoCashouts(currentMultiplier float64) error {
	// 查找需要自动止盈的下注
	var bets []model.Bet
	err := s.db.Where("status = 0 AND auto_cashout > 0 AND auto_cashout <= ?", currentMultiplier).Find(&bets).Error
	if err != nil {
		return err
	}

	// 处理每个自动止盈
	for _, bet := range bets {
		payout := bet.Amount * currentMultiplier
		
		// 更新下注状态
		if err := s.UpdateBetCashout(bet.BetID, currentMultiplier, payout); err != nil {
			continue
		}

		// 增加用户余额
		if err := s.AddUserBalance(bet.UserID, payout); err != nil {
			continue
		}

		// 更新用户统计
		s.UpdateUserStats(bet.UserID, bet.Amount, payout, currentMultiplier)
	}

	return nil
}

// ProcessCrashedBets 处理崩盘的下注
func (s *GameService) ProcessCrashedBets() error {
	// 将所有进行中的下注标记为崩盘
	return s.db.Model(&model.Bet{}).Where("status = 0").Update("status", 2).Error
}
