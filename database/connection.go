package database

import (
	"fmt"
	"sync"

	"github.com/chmenegatti/nsxt-vs/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

func GetDatabaseConnection(edge string, cfg config.Config, logger *zap.Logger) (*gorm.DB, error) {
	var err error
	dbOnce.Do(
		func() {
			conf := cfg.Server[edge]
			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port,
				conf.DBName,
			)
			dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				logger.Error("Failed to connect to database", zap.Error(err))
				return
			}
		},
	)
	return dbInstance, err
}
