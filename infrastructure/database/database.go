package database

import (
	"context"
	"embed"

	_ "github.com/lib/pq"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
	"gorm.io/driver/postgres"
)

var (
	//go:embed migration/postgres/*.sql
	embedMigrations embed.FS
	dbConnections   DatabaseConnections
)

// DatabaseConnections represent all available database connections
type DatabaseConfigs struct {
	Postgres PgConfigs
	MySQL    MySQLConfigs
	Mongo    MongoConfigs
}

// DatabaseConnections represent all available database connections
type DatabaseConnections struct {
	Postgres PgConnections
	MySQL    MySQLConnections
	Mongo    MongoConnections
}

const (
	SamplePostgresDBName = "SamplePostgresDB"
	SampleMySQLDBName    = "sampleMySQLDB"
	SampleMongoDBName    = "sampleMongoDB"
)

func GetConfigs() DatabaseConfigs {

	// TODO: get Database Config implementation simpler through config file
	var (
		pgConfigs    = make(map[string]PostgresConfig)
		mysqlConfigs = make(map[string]MySQLConfig)
		mongoConfigs = make(map[string]MongoConfig)
	)

	pgConfig := PostgresConfig{
		Name: SamplePostgresDBName,
		Master: postgres.Config{
			DSN:                  secret.GetPostgresMasterDSN(),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		},
		Slave: postgres.Config{
			DSN:                  secret.GetPostgresSlaveDSN(),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		},
	}
	pgConfigs[pgConfig.Name] = pgConfig

	mongoConfig := MongoConfig{
		Name:     SampleMongoDBName,
		Timeout:  secret.MONGO_TIMEOUT,
		DBname:   secret.MONGO_DB,
		Username: secret.MONGO_USER,
		Password: secret.MONGO_PASSWORD,
		Host:     secret.MONGO_HOST,
		Port:     secret.MONGO_PORT,
	}
	mongoConfigs[mongoConfig.Name] = mongoConfig

	return DatabaseConfigs{
		Postgres: pgConfigs,
		MySQL:    mysqlConfigs,
		Mongo:    mongoConfigs,
	}
}

// CreateDBConnection creates all database connections used in this app (postgre, mysql, mongo, etc)
//
// It will close the existing DB connection before opening a new DB connection
func CreateDBConnection(ctx context.Context, dbConfigs DatabaseConfigs) DatabaseConnections {

	closeDBConnections(ctx, dbConnections)

	return openDBConnections(ctx, dbConfigs)

}

// openDBConnections opend database connection for each database
func openDBConnections(ctx context.Context, dbConfigs DatabaseConfigs) DatabaseConnections {
	var (
		pgConnections    PgConnections
		mysqlConnections MySQLConnections
		mongoConnections MongoConnections
	)

	if len(dbConfigs.Postgres) > 0 {
		pgConnections = openPostgresConnections(dbConfigs.Postgres)
	}

	if len(dbConfigs.MySQL) > 0 {
		mysqlConnections = openMySQLConnections(dbConfigs.MySQL)
	}

	if len(dbConfigs.Mongo) > 0 {
		mongoConnections = openMongoConnections(ctx, dbConfigs.Mongo)
	}

	return DatabaseConnections{
		Postgres: pgConnections,
		MySQL:    mysqlConnections,
		Mongo:    mongoConnections,
	}
}

// closeDBConnections close the existing connections if open
func closeDBConnections(ctx context.Context, dbConnections DatabaseConnections) {
	for _, dbPostgres := range dbConnections.Postgres {
		if dbPostgres.Master != nil {
			closePostgres(dbPostgres.Master)
		}

		if dbPostgres.Slave != nil {
			closePostgres(dbPostgres.Slave)
		}
	}

	for _, dbMySQL := range dbConnections.MySQL {
		if dbMySQL.Master != nil {
			closeMySQL(dbMySQL.Master)
		}

		if dbMySQL.Slave != nil {
			closeMySQL(dbMySQL.Slave)
		}
	}

	for _, dbMongo := range dbConnections.Mongo {
		if dbMongo.DB != nil {
			closeMongo(ctx, dbMongo.DB)
		}
	}
}

// openPostgresConnections open Postgres connections for Master and Slave database
func openPostgresConnections(pgConfigs PgConfigs) PgConnections {
	pgConnections := make(map[string]PgConnection)

	// TODO: update open DB with concurrency
	for _, pgConfig := range pgConfigs {
		pgConnections[pgConfig.Name] = PgConnection{
			Name:   pgConfig.Name,
			Master: openPostgres(pgConfig.Master),
			Slave:  openPostgres(pgConfig.Slave),
		}
	}
	return pgConnections
}

// openMySQLConnections open Postgres connections for Master and Slave database
func openMySQLConnections(mysqlConfigs MySQLConfigs) MySQLConnections {
	mysqlConnections := make(map[string]MySQLConnection)

	// TODO: update open DB with concurrency
	for _, mysqlConfig := range mysqlConfigs {
		mysqlConnections[mysqlConfig.Name] = MySQLConnection{
			Name:   mysqlConfig.Name,
			Master: openMySQL(mysqlConfig.Master),
			Slave:  openMySQL(mysqlConfig.Slave),
		}
	}
	return mysqlConnections
}

// openMySQLConnections open Postgres connections for Master and Slave database
func openMongoConnections(ctx context.Context, mongoConfigs MongoConfigs) MongoConnections {
	mongoConnections := make(map[string]MongoConnection)

	// TODO: update open DB with concurrency
	for _, mongoConfig := range mongoConfigs {
		mongoConnections[mongoConfig.Name] = MongoConnection{
			Name: mongoConfig.Name,
			DB:   openMongo(ctx, mongoConfig),
		}
	}
	return mongoConnections
}
