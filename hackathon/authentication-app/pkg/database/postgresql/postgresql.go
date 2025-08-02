package database

import (
	"authentication-app/config"
	"context"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(c *config.Config) (*gorm.DB, error) {
	dbHost := c.DBHost
	dbUser := c.DBUser
	dbPassword := c.DBPassword
	dbName := c.DBName
	dbPort := c.DBPort
	dbSSLMode := c.DBSSLMode
	dbSchema := c.DBSchema
	hosts := strings.Split(dbHost, ",")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s&target_session_attrs=read-write",
		dbUser, dbPassword, strings.Join(hosts, ","), dbPort, dbName, dbSSLMode, dbSchema)

	logMode := logger.Silent
	if c.ServerDebug {
		logMode = logger.Info
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        dsn,
	}), &gorm.Config{
		Logger:               logger.Default.LogMode(logMode),
		DisableAutomaticPing: false,
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(c.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.DBConnMaxLifeTime) * time.Second)

	return db, nil
}

func IsPostgreSQLReady(ctx context.Context, db *gorm.DB) (isReady bool) {
	d, err := db.DB()
	if err != nil {
		return
	}
	return d.PingContext(ctx) == nil
}
