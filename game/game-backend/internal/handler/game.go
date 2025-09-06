package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"game-backend/internal/middleware"
	"game-backend/internal/service"
	"game-backend/internal/websocket"
)

// GameHandler 游戏处理器
type GameHandler struct {
	gameService *service.GameService
	wsHub       *websocket.Hub
}

// NewGameHandler 创建游戏处理器
func NewGameHandler(gameService *service.GameService, wsHub *websocket.Hub) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		wsHub:       wsHub,
	}
}

// BetRequest 下注请求结构
type BetRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=1"`
	AutoCashout float64 `json:"auto_cashout" binding:"omitempty,min=1.01"`
}

// CashoutRequest 止盈请求结构
type CashoutRequest struct {
	BetID string `json:"bet_id" binding:"required"`
}

// GetGameStatus 获取游戏状态
func (h *GameHandler) GetGameStatus(c *gin.Context) {
	gameState := h.wsHub.GetGameState()
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"game_id":            gameState.GameID,
			"status":             gameState.Status,
			"current_multiplier": gameState.CurrentMultiplier,
			"players_count":      gameState.PlayersCount,
			"next_round_in":      gameState.NextRoundIn,
			"server_time":        gameState.LastUpdate,
		},
	})
}

// PlaceBet 下注
func (h *GameHandler) PlaceBet(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	var req BetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证下注金额
	if req.Amount < h.gameService.GetMinBetAmount() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "下注金额不能小于最小限制",
		})
		return
	}

	if req.Amount > h.gameService.GetMaxBetAmount() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "下注金额不能大于最大限制",
		})
		return
	}

	// 检查用户余额
	user, err := h.gameService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	if user.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "余额不足",
		})
		return
	}

	// 创建下注记录
	bet, err := h.gameService.CreateBet(userID, req.Amount, req.AutoCashout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "下注失败: " + err.Error(),
		})
		return
	}

	// 扣除用户余额
	err = h.gameService.DeductUserBalance(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "扣除余额失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "下注成功",
		"data": gin.H{
			"bet_id":       bet.BetID,
			"amount":       bet.Amount,
			"auto_cashout": bet.AutoCashout,
			"status":       bet.Status,
		},
	})
}

// Cashout 止盈
func (h *GameHandler) Cashout(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	var req CashoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取下注记录
	bet, err := h.gameService.GetBetByID(req.BetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "下注记录不存在",
		})
		return
	}

	// 检查下注是否属于当前用户
	if bet.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权操作此下注",
		})
		return
	}

	// 检查下注状态
	if bet.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "下注已处理",
		})
		return
	}

	// 获取当前游戏状态
	gameState := h.wsHub.GetGameState()
	if gameState.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "游戏未进行中",
		})
		return
	}

	// 计算赔付
	payout := bet.Amount * gameState.CurrentMultiplier
	profit := payout - bet.Amount

	// 更新下注记录
	err = h.gameService.UpdateBetCashout(req.BetID, gameState.CurrentMultiplier, payout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "止盈失败: " + err.Error(),
		})
		return
	}

	// 增加用户余额
	err = h.gameService.AddUserBalance(userID, payout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "增加余额失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "止盈成功",
		"data": gin.H{
			"bet_id":     bet.BetID,
			"multiplier": gameState.CurrentMultiplier,
			"payout":     payout,
			"profit":     profit,
		},
	})
}

// GetBetHistory 获取下注历史
func (h *GameHandler) GetBetHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取下注历史
	bets, total, err := h.gameService.GetUserBetHistory(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取历史失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"bets":      bets,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetGameHistory 获取游戏历史
func (h *GameHandler) GetGameHistory(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	// 获取游戏历史
	games, total, err := h.gameService.GetGameHistory(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取历史失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"games":     games,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetLeaderboard 获取排行榜
func (h *GameHandler) GetLeaderboard(c *gin.Context) {
	// 获取排行榜
	leaderboard, err := h.gameService.GetLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取排行榜失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    leaderboard,
	})
}

// GetUserStats 获取用户统计
func (h *GameHandler) GetUserStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	// 获取用户统计
	stats, err := h.gameService.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    stats,
	})
}
