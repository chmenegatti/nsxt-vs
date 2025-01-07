package database

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LoadBalance struct {
	VIPPort            string `gorm:"column:vip_port"`
	AddressLoadBalance string `gorm:"column:address_load_balance"`
	NSXTServerID       string `gorm:"column:nsxt_virtual_server_id"`
}

func FetchLoadBalances(db *gorm.DB, logger *zap.Logger) ([]LoadBalance, error) {
	var records []LoadBalance
	if err := db.Table("load_balances").Select("vip_port, address_load_balance, nsxt_virtual_server_id").Find(&records).Error; err != nil {
		logger.Error("Failed to fetch data", zap.Error(err))
		return nil, err
	}
	return records, nil
}
