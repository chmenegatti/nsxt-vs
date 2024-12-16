package database

import (
	"database/sql"
	"fmt"

	"github.com/chmenegatti/nsxt-vs/config"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseManager struct {
	DB *sql.DB
}

func NewDatabaseManager(cfg config.DatabaseConfig) (*DatabaseManager, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("database is not reachable: %w", err)
	}

	return &DatabaseManager{DB: db}, nil
}

func (dm *DatabaseManager) QueryLoadBalances() ([][3]string, error) {
	const query = "SELECT vip_port, address_load_balance, nsxt_virtual_server_id FROM load_balances"
	rows, err := dm.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results [][3]string
	for rows.Next() {
		var vipPort, addressLoadBalance, nsxtVirtualServerID string
		if err := rows.Scan(&vipPort, &addressLoadBalance, &nsxtVirtualServerID); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		id := nsxtVirtualServerID
		displayName := fmt.Sprintf("%s-%s", addressLoadBalance, vipPort)
		results = append(results, [3]string{id, displayName, vipPort})
	}

	return results, nil
}

func (dm *DatabaseManager) Close() error {
	if dm.DB != nil {
		return dm.DB.Close()
	}
	return nil
}
