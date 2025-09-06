package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Game     GameConfig     `mapstructure:"game"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"` // debug, release, test
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"` // 小时
	Issuer     string `mapstructure:"issuer"`
}

// GameConfig 游戏配置
type GameConfig struct {
	MinBetAmount      float64 `mapstructure:"min_bet_amount"`
	MaxBetAmount      float64 `mapstructure:"max_bet_amount"`
	MinMultiplier     float64 `mapstructure:"min_multiplier"`
	MaxMultiplier     float64 `mapstructure:"max_multiplier"`
	RoundDuration     int     `mapstructure:"round_duration"`     // 秒
	BettingDuration   int     `mapstructure:"betting_duration"`  // 秒
	WaitingDuration   int     `mapstructure:"waiting_duration"`   // 秒
	UpdateInterval    int     `mapstructure:"update_interval"`    // 毫秒
	MaxPlayersPerGame int     `mapstructure:"max_players_per_game"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("读取配置文件失败: %v", err)
		return err
	}

	// 解析配置
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Printf("解析配置失败: %v", err)
		return err
	}

	// 验证配置
	if err := validateConfig(); err != nil {
		log.Printf("配置验证失败: %v", err)
		return err
	}

	log.Printf("配置加载成功: %s", configPath)
	return nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	// 数据库默认配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "crash_game")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle", 10)
	viper.SetDefault("database.max_open", 100)

	// Redis默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// JWT默认配置
	viper.SetDefault("jwt.secret", "crash-game-secret-key")
	viper.SetDefault("jwt.expire_time", 24)
	viper.SetDefault("jwt.issuer", "crash-game")

	// 游戏默认配置
	viper.SetDefault("game.min_bet_amount", 1.0)
	viper.SetDefault("game.max_bet_amount", 1000.0)
	viper.SetDefault("game.min_multiplier", 1.01)
	viper.SetDefault("game.max_multiplier", 1000.0)
	viper.SetDefault("game.round_duration", 30)
	viper.SetDefault("game.betting_duration", 15)
	viper.SetDefault("game.waiting_duration", 10)
	viper.SetDefault("game.update_interval", 100)
	viper.SetDefault("game.max_players_per_game", 1000)

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 7)
}

// validateConfig 验证配置
func validateConfig() error {
	if AppConfig.Server.Port <= 0 || AppConfig.Server.Port > 65535 {
		return fmt.Errorf("服务器端口无效: %d", AppConfig.Server.Port)
	}

	if AppConfig.Database.Port <= 0 || AppConfig.Database.Port > 65535 {
		return fmt.Errorf("数据库端口无效: %d", AppConfig.Database.Port)
	}

	if AppConfig.Redis.Port <= 0 || AppConfig.Redis.Port > 65535 {
		return fmt.Errorf("Redis端口无效: %d", AppConfig.Redis.Port)
	}

	if AppConfig.Game.MinBetAmount <= 0 {
		return fmt.Errorf("最小下注金额必须大于0")
	}

	if AppConfig.Game.MaxBetAmount <= AppConfig.Game.MinBetAmount {
		return fmt.Errorf("最大下注金额必须大于最小下注金额")
	}

	if AppConfig.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
}

// GetRedisAddr 获取Redis地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetServerAddr 获取服务器地址
func (c *ServerConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDebug 是否为调试模式
func (c *ServerConfig) IsDebug() bool {
	return c.Mode == "debug"
}

// GetEnvValue 获取环境变量值，如果不存在则返回默认值
func GetEnvValue(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt 获取环境变量整数值，如果不存在则返回默认值
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvFloat 获取环境变量浮点数值，如果不存在则返回默认值
func GetEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
