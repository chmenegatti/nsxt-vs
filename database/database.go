package database

import (
	"database/sql"
	"fmt"

	"github.com/chmenegatti/nsxt-vs/config"
	"github.com/chmenegatti/nsxt-vs/utils"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseManager struct {
	db *sql.DB
}

func NewDatabaseManager(cfg config.DatabaseConfig) (*DatabaseManager, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}
	return &DatabaseManager{db: db}, nil
}

func (dm *DatabaseManager) QueryLoadBalances() ([][3]string, error) {
	query := `SELECT vip_port, address_load_balance, nsxt_virtual_server_id FROM load_balances`
	rows, err := dm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results [][3]string
	for rows.Next() {
		var vipPort, addressLoadBalance, nsxtVirtualServerID string
		if err := rows.Scan(&vipPort, &addressLoadBalance, &nsxtVirtualServerID); err != nil {
			return nil, err
		}
		displayName := fmt.Sprintf("%s-%s", addressLoadBalance, vipPort)
		results = append(results, [3]string{nsxtVirtualServerID, displayName, vipPort})
	}

	utils.SortLoadBalancesByIP(results)
	return results, nil
}

func (dm *DatabaseManager) Close() {
	if dm.db != nil {
		dm.db.Close()
	}
}
