// Package secrets are values that supposed to be secrets and private and not public consuming
// for example: keys, passwords, ip addresses
//
// To differentiate with config, we place secrets on a different package, so it will be called from outside like this:
// `secrets.POSTGRES_HOST` for instance
// You just need to add more env variables here
package secret

import (
	"fmt"
	"strings"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/scheduler"
	"github.com/spf13/viper"
)

var POSTGRES_HOST_MASTER string
var POSTGRES_PORT_MASTER int
var POSTGRES_USERNAME_MASTER string
var POSTGRES_PASSWORD_MASTER string
var POSTGRES_DATABASE_MASTER string
var POSTGRES_SSL_MODE_MASTER string

var POSTGRES_HOST_SLAVE string
var POSTGRES_PORT_SLAVE int
var POSTGRES_USERNAME_SLAVE string
var POSTGRES_PASSWORD_SLAVE string
var POSTGRES_DATABASE_SLAVE string
var POSTGRES_SSL_MODE_SLAVE string

var MYSQL_HOST_MASTER string
var MYSQL_PORT_MASTER int
var MYSQL_USERNAME_MASTER string
var MYSQL_PASSWORD_MASTER string
var MYSQL_DATABASE_MASTER string

var MYSQL_HOST_SLAVE string
var MYSQL_PORT_SLAVE int
var MYSQL_USERNAME_SLAVE string
var MYSQL_PASSWORD_SLAVE string
var MYSQL_DATABASE_SLAVE string

var MONGO_HOST string
var MONGO_PORT string
var MONGO_DB string
var MONGO_USER string
var MONGO_PASSWORD string
var MONGO_TIMEOUT int

var JWT_METHOD string
var JWT_SECRET string
var JWT_LIFESPAN_AUTHTOKEN int
var JWT_LIFESPAN_REFRESHTOKEN int
var JWT_DOMAIN string
var JWT_AUDIENCE []string
var JWKS_URL string
var JWKS_REFRESH int
var JWKS_TTL int

var REDIS_HOST string
var REDIS_PORT int
var REDIS_PASSWORD string

var RUNDECK_URL string
var RUNDECK_API_TOKEN string
var RUNDECK_PROJECT string
var RUNDECK_JOB_ID string
var RUNDECK_JOB_BEARER_TOKEN string

var NEWRELIC_APP_NAME string
var NEWRELIC_LICENSE string

var HYDRA_PUBLIC_URL string
var HYDRA_CLIENT_ID string
var HYDRA_CLIENT_SECRET string
var HYDRA_AUDIENCE string

var RABBITMQ_URI string

// Reload reload secret either from file or from system's ENV
// see: infrastructure/configuration/setup.go
func Reload() {
	JWT_METHOD = viper.GetString("JWT_METHOD")
	JWT_SECRET = viper.GetString("JWT_SECRET")
	JWT_LIFESPAN_AUTHTOKEN = viper.GetInt("JWT_LIFESPAN_AUTHTOKEN")
	JWT_LIFESPAN_REFRESHTOKEN = viper.GetInt("JWT_LIFESPAN_REFRESHTOKEN")
	JWT_DOMAIN = viper.GetString("JWT_DOMAIN")
	JWT_AUDIENCE = strings.Split(viper.GetString("JWT_AUDIENCE"), ",")
	JWKS_URL = viper.GetString("JWKS_URL")
	JWKS_REFRESH = viper.GetInt("JWKS_REFRESH")
	JWKS_TTL = viper.GetInt("JWKS_TTL")
	REDIS_HOST = viper.GetString("REDIS_HOST")
	REDIS_PORT = viper.GetInt("REDIS_PORT")
	REDIS_PASSWORD = viper.GetString("REDIS_PASSWORD")
	RUNDECK_URL = viper.GetString("RUNDECK_URL")
	RUNDECK_API_TOKEN = viper.GetString("RUNDECK_API_TOKEN")
	RUNDECK_PROJECT = viper.GetString("RUNDECK_PROJECT")
	RUNDECK_JOB_ID = viper.GetString("RUNDECK_JOB_ID")
	RUNDECK_JOB_BEARER_TOKEN = viper.GetString("RUNDECK_JOB_BEARER_TOKEN")
	NEWRELIC_APP_NAME = viper.GetString("NEWRELIC_APP_NAME")
	NEWRELIC_LICENSE = viper.GetString("NEWRELIC_LICENSE")
	HYDRA_PUBLIC_URL = viper.GetString("HYDRA_PUBLIC_URL")
	HYDRA_CLIENT_ID = viper.GetString("HYDRA_CLIENT_ID")
	HYDRA_CLIENT_SECRET = viper.GetString("HYDRA_CLIENT_SECRET")
	HYDRA_AUDIENCE = viper.GetString("HYDRA_AUDIENCE")
	RABBITMQ_URI = viper.GetString("RABBITMQ_URI")

	POSTGRES_HOST_MASTER = viper.GetString("POSTGRES_HOST_MASTER")
	POSTGRES_PORT_MASTER = viper.GetInt("POSTGRES_PORT_MASTER")
	POSTGRES_USERNAME_MASTER = viper.GetString("POSTGRES_USERNAME_MASTER")
	POSTGRES_PASSWORD_MASTER = viper.GetString("POSTGRES_PASSWORD_MASTER")
	POSTGRES_DATABASE_MASTER = viper.GetString("POSTGRES_DATABASE_MASTER")
	POSTGRES_SSL_MODE_MASTER = viper.GetString("POSTGRES_SSL_MODE_MASTER")

	POSTGRES_HOST_SLAVE = viper.GetString("POSTGRES_HOST_SLAVE")
	POSTGRES_PORT_SLAVE = viper.GetInt("POSTGRES_PORT_SLAVE")
	POSTGRES_USERNAME_SLAVE = viper.GetString("POSTGRES_USERNAME_SLAVE")
	POSTGRES_PASSWORD_SLAVE = viper.GetString("POSTGRES_PASSWORD_SLAVE")
	POSTGRES_DATABASE_SLAVE = viper.GetString("POSTGRES_DATABASE_SLAVE")
	POSTGRES_SSL_MODE_SLAVE = viper.GetString("POSTGRES_SSL_MODE_SLAVE")

	MYSQL_HOST_MASTER = viper.GetString("MYSQL_HOST_MASTER")
	MYSQL_PORT_MASTER = viper.GetInt("MYSQL_PORT_MASTER")
	MYSQL_USERNAME_MASTER = viper.GetString("MYSQL_USERNAME_MASTER")
	MYSQL_PASSWORD_MASTER = viper.GetString("MYSQL_PASSWORD_MASTER")
	MYSQL_DATABASE_MASTER = viper.GetString("MYSQL_DATABASE_MASTER")

	MYSQL_HOST_SLAVE = viper.GetString("MYSQL_HOST_SLAVE")
	MYSQL_PORT_SLAVE = viper.GetInt("MYSQL_PORT_SLAVE")
	MYSQL_USERNAME_SLAVE = viper.GetString("MYSQL_USERNAME_SLAVE")
	MYSQL_PASSWORD_SLAVE = viper.GetString("MYSQL_PASSWORD_SLAVE")
	MYSQL_DATABASE_SLAVE = viper.GetString("MYSQL_DATABASE_SLAVE")

	MONGO_HOST = viper.GetString("MONGO_HOST")
	MONGO_PORT = viper.GetString("MONGO_PORT")
	MONGO_DB = viper.GetString("MONGO_DB")
	MONGO_USER = viper.GetString("MONGO_USER")
	MONGO_PASSWORD = viper.GetString("MONGO_PASSWORD")
	MONGO_TIMEOUT = viper.GetInt("MONGO_TIMEOUT")

}

// dsnConfig represent data source name (DSN) configuration
type dsnConfig struct {
	host     string
	user     string
	password string
	db       string
	port     int
	ssl      string
}

// GetRedisDSN get Redis data source name
func GetRedisDSN() string {
	var (
		host string = REDIS_HOST
		port int    = REDIS_PORT
	)

	return fmt.Sprintf("%s:%v", host, port)
}

// GetPostgresMasterDSN get Postgres's master data source name
func GetPostgresMasterDSN() string {
	return writePostgreDSNString(dsnConfig{
		host:     POSTGRES_HOST_MASTER,
		user:     POSTGRES_USERNAME_MASTER,
		password: POSTGRES_PASSWORD_MASTER,
		db:       POSTGRES_DATABASE_MASTER,
		port:     POSTGRES_PORT_MASTER,
		ssl:      POSTGRES_SSL_MODE_MASTER,
	})
}

// GetPostgresMasterDSN get Postgres's slave data source name
func GetPostgresSlaveDSN() string {
	return writePostgreDSNString(dsnConfig{
		host:     POSTGRES_HOST_SLAVE,
		user:     POSTGRES_USERNAME_SLAVE,
		password: POSTGRES_PASSWORD_SLAVE,
		db:       POSTGRES_DATABASE_SLAVE,
		port:     POSTGRES_PORT_SLAVE,
		ssl:      POSTGRES_SSL_MODE_SLAVE,
	})
}

// GetMySQLMasterDSN get MySQL's master data source name
func GetMySQLMasterDSN() string {
	return writeMySQLDSNString(dsnConfig{
		host:     MYSQL_HOST_MASTER,
		user:     MYSQL_USERNAME_MASTER,
		password: MYSQL_PASSWORD_MASTER,
		db:       MYSQL_DATABASE_MASTER,
		port:     MYSQL_PORT_MASTER,
	})
}

// GetMySQLSlaveDSN get MySQL's slave data source name
func GetMySQLSlaveDSN() string {
	return writeMySQLDSNString(dsnConfig{
		host:     MYSQL_HOST_SLAVE,
		user:     MYSQL_USERNAME_SLAVE,
		password: MYSQL_PASSWORD_SLAVE,
		db:       MYSQL_DATABASE_SLAVE,
		port:     MYSQL_PORT_SLAVE,
	})
}

// GetSchedulerConfig get scheduler config config from the config vars
func GetSchedulerConfig() *scheduler.SchedulerConfig {
	return &scheduler.SchedulerConfig{
		RundeckUrl:      RUNDECK_URL,
		RundeckAPIToken: RUNDECK_API_TOKEN,
		RundeckProject:  RUNDECK_PROJECT,
	}
}

// writePostgreDSNString write Postgres DSN format string
func writePostgreDSNString(dsn dsnConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s", dsn.host, dsn.user, dsn.password, dsn.db, dsn.port, dsn.ssl)
}

// writeMySQLDSNString write MySQL DSN format string
func writeMySQLDSNString(dsn dsnConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", dsn.user, dsn.password, dsn.host, dsn.port, dsn.db)
}
