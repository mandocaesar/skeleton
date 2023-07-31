package database

import (
	"context"
	"database/sql"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// MySQLConfigs represent collections of MySQL database configurations
type MySQLConfigs map[string]MySQLConfig

// MySQLConnections represent collections of MySQL database connections
type MySQLConnections map[string]MySQLConnection

// MySQLConfig represent MySQL database configuration
type MySQLConfig struct {
	Name   string
	Master mysql.Config
	Slave  mysql.Config
}

// MySQLConnection represent MySQL database connections
type MySQLConnection struct {
	Name   string
	Master *gorm.DB
	Slave  *gorm.DB
}

const (
	DB_DIALECT_MYSQL    = "mysql"
	MIGRATION_DIR_MYSQL = "migration/mysql"
)

func openMySQL(mysqlConfig mysql.Config) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "gorm.Open() got error on connecting to database - database.openMySQL()")
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithoutQueryVariables())); err != nil {
		log.StdFatal(context.Background(), nil, err, "db.Use(tracing.NewPlugin(tracing.WithoutQueryVariables())) got error  - database.openMySQL()")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "db.DB() got error on opening database - database.openMySQL()")
	}

	if err = sqlDB.Ping(); err != nil {
		log.StdFatal(context.Background(), nil, err, "sqlDB.Ping() got error on ping database - database.openMySQL()")
	}

	if config.DB_MYSQL_AUTO_MIGRATE {
		runMySQLMigration(sqlDB)
	}

	log.Info("successfully connected to mysql database")

	return db
}

func closeMySQL(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "conn.DB()Error occurred while closing a DB connection - database.closeMySQL()")
	}

	sqlDB.Close()
}

func runMySQLMigration(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(DB_DIALECT_MYSQL); err != nil {
		log.StdFatal(context.Background(), nil, err, "goose.SetDialect() got error - runMySQLMigration")
	}

	if err := goose.Up(db, MIGRATION_DIR_MYSQL); err != nil {
		log.StdFatal(context.Background(), nil, err, " goose.Up() got error - runMySQLMigration")
	}

	log.Info("mysql database migrations applied")
}
