package util

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fliqt/config"
)

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.GetDBDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	} else {
		db.Logger = logger.Default.LogMode(logger.Silent)
	}

	sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
	sqlDB.SetMaxOpenConns(cfg.DBMaxConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBMaxLifeTime))

	return db, nil
}
