-- 创建数据库和用户
CREATE DATABASE IF NOT EXISTS crash_game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（如果不存在）
CREATE USER IF NOT EXISTS 'crash_user'@'%' IDENTIFIED BY 'crash_password';
GRANT ALL PRIVILEGES ON crash_game.* TO 'crash_user'@'%';
FLUSH PRIVILEGES;

-- 使用数据库
USE crash_game;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    balance DECIMAL(15,2) DEFAULT 0.00,
    avatar VARCHAR(255),
    status TINYINT DEFAULT 1 COMMENT '1:正常 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户统计表
CREATE TABLE IF NOT EXISTS user_stats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    total_bets BIGINT DEFAULT 0,
    total_winnings DECIMAL(15,2) DEFAULT 0.00,
    biggest_multiplier DECIMAL(10,2) DEFAULT 0.00,
    games_played BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_id (user_id),
    INDEX idx_total_winnings (total_winnings),
    INDEX idx_biggest_multiplier (biggest_multiplier)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token VARCHAR(500) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_token (token),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建游戏表
CREATE TABLE IF NOT EXISTS games (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    game_id VARCHAR(50) NOT NULL UNIQUE,
    status TINYINT DEFAULT 0 COMMENT '0:等待 1:进行中 2:已结束',
    round_id VARCHAR(50),
    multiplier DECIMAL(10,2) DEFAULT 0.00,
    players_count INT DEFAULT 0,
    total_bets DECIMAL(15,2) DEFAULT 0.00,
    total_payout DECIMAL(15,2) DEFAULT 0.00,
    start_time TIMESTAMP NULL,
    end_time TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_game_id (game_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建下注表
CREATE TABLE IF NOT EXISTS bets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    bet_id VARCHAR(50) NOT NULL UNIQUE,
    user_id BIGINT UNSIGNED NOT NULL,
    game_id VARCHAR(50) NOT NULL,
    round_id VARCHAR(50),
    amount DECIMAL(15,2) NOT NULL,
    auto_cashout DECIMAL(10,2) DEFAULT 0.00,
    multiplier DECIMAL(10,2) DEFAULT 0.00,
    payout DECIMAL(15,2) DEFAULT 0.00,
    status TINYINT DEFAULT 0 COMMENT '0:进行中 1:已止盈 2:已崩盘',
    cashout_time TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_bet_id (bet_id),
    INDEX idx_user_id (user_id),
    INDEX idx_game_id (game_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建游戏历史表
CREATE TABLE IF NOT EXISTS game_history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    round_id VARCHAR(50) NOT NULL UNIQUE,
    game_id VARCHAR(50) NOT NULL,
    final_multiplier DECIMAL(10,2) NOT NULL,
    players_count INT DEFAULT 0,
    total_bets DECIMAL(15,2) DEFAULT 0.00,
    total_payout DECIMAL(15,2) DEFAULT 0.00,
    winners_count INT DEFAULT 0,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_round_id (round_id),
    INDEX idx_game_id (game_id),
    INDEX idx_final_multiplier (final_multiplier),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建排行榜表
CREATE TABLE IF NOT EXISTS leaderboard (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    username VARCHAR(50) NOT NULL,
    total_winnings DECIMAL(15,2) DEFAULT 0.00,
    biggest_multiplier DECIMAL(10,2) DEFAULT 0.00,
    rank INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_total_winnings (total_winnings),
    INDEX idx_rank (rank)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入测试用户
INSERT IGNORE INTO users (username, password, email, balance, status) VALUES
('testuser1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'test1@example.com', 1000.00, 1),
('testuser2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'test2@example.com', 1000.00, 1);

-- 插入测试用户统计
INSERT IGNORE INTO user_stats (user_id, total_bets, total_winnings, biggest_multiplier, games_played) VALUES
(1, 0, 0.00, 0.00, 0),
(2, 0, 0.00, 0.00, 0);
