package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"game-backend/config"
	"game-backend/internal/model"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL数据库连接
func InitMySQL() error {
	// 构建DSN
	dsn := config.AppConfig.Database.GetDSN()
	
	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取底层sql.DB对象
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdle)
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	DB = db
	log.Println("MySQL数据库连接成功")
	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 自动迁移所有模型
	err := DB.AutoMigrate(
		&model.User{},
		&model.UserStats{},
		&model.UserSession{},
		&model.Game{},
		&model.Bet{},
		&model.GameHistory{},
		&model.Leaderboard{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Println("数据库表结构迁移完成")
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
