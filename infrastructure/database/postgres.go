package database

import (
	"context"
	"database/sql"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

const (
	DB_DIALECT_POSTGRES    = "postgres"
	MIGRATION_DIR_POSTGRES = "migration/postgres"
)

// PostgresConfigs represent collection of Postgres database configuration
type PgConfigs map[string]PostgresConfig

// PgConnections represent collection of Postgres database connections
type PgConnections map[string]PgConnection

// PostgresConfig represent Postgres database configuration
type PostgresConfig struct {
	Name   string
	Master postgres.Config
	Slave  postgres.Config
}

// PgConnection represent Postgres database connections
type PgConnection struct {
	Name   string
	Master *gorm.DB
	Slave  *gorm.DB
}

func openPostgres(pgConfig postgres.Config) *gorm.DB {
	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "gorm.Open() got error on connecting to database - database.openPostgres()")
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithoutQueryVariables())); err != nil {
		log.StdFatal(context.Background(), nil, err, "db.Use(tracing.NewPlugin(tracing.WithoutQueryVariables())) got error  - database.openPostgres()")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "db.DB() got error on opening database - database.openPostgres()")
	}

	if err = sqlDB.Ping(); err != nil {
		log.StdFatal(context.Background(), nil, err, "sqlDB.Ping() got error on ping database - database.openPostgres()")
	}

	if config.DB_POSTGRES_AUTO_MIGRATE {
		runPostgresMigration(sqlDB)
	}

	log.Info("successfully connected to postgres database")

	return db
}

func closePostgres(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "conn.DB()Error occurred while closing a DB connection - database.closePostgres()")
	}

	sqlDB.Close()
}

func runPostgresMigration(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(DB_DIALECT_POSTGRES); err != nil {
		log.StdFatal(context.Background(), nil, err, "goose.SetDialect() got error - runPostgresMigration()")
	}

	if err := goose.Up(db, MIGRATION_DIR_POSTGRES); err != nil {
		log.StdFatal(context.Background(), nil, err, " goose.Up() got error - runPostgresMigration()")
	}

	log.Info("postgres database migrations applied")
}
