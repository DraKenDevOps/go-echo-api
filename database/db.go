package database

import (
	"database/sql"
	"fmt"
	"go-echo-api/config"
	"go-echo-api/zplogger"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config, logger *zplogger.Logger) *sql.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
	}

	return db
}
