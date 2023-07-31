package config

import (
	"time"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/cache"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
	rmqConfig "github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/common/config"
	flagConfig "github.com/machtwatch/catalystdk/go/flag/config"
	"github.com/machtwatch/catalystdk/go/log"
	traceConfig "github.com/machtwatch/catalystdk/go/trace/config"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

// You just need to add more env variables here
var APP_NAME string
var APP_PORT int
var APP_ENV string
var APP_URL string
var CACHE_WHITELIST_TOKENS string
var CONTEXT_TIMEOUT time.Duration
var DB_POSTGRES_AUTO_MIGRATE bool
var DB_MYSQL_AUTO_MIGRATE bool
var DB_TIMEZONE string
var DEBUG bool
var BIFROST_URI string
var BIFROST_TENANT_CODE string
var FRONTEND_URL string
var IMAGE_MAX_SIZE_IN_MB int64
var IMAGE_MAX_WIDTH int
var IMAGE_MAX_HEIGHT int
var IMAGE_MIN_WIDTH int
var IMAGE_MIN_HEIGHT int
var HTTP_REQUEST_RETRY_COUNT int
var HTTP_REQUEST_RETRY_WAIT_TIME_MS int
var HTTP_REQUEST_MAX_WAIT_TIME_MS int
var JWT_METHOD string
var KRATOS_URL string
var LOG_CALLER bool
var LOG_LEVEL string
var LOG_USE_JSON bool
var NOTIFICATION_SERVICE_BASE_URL string
var NOTIFICATION_SERVICE_SENDER_EMAIL string
var NOTIFICATION_SERVICE_SENDER_NAME string
var NR_HOST string
var NR_METRIC_PREFIX string
var NR_LICENSE_KEY string
var RABBITMQ_CONFIG_DURABLE bool
var RABBITMQ_CONFIG_EXCLUSIVE bool
var RABBITMQ_CONFIG_AUTO_DELETED bool
var RABBITMQ_CONFIG_INTERNAL bool
var RABBITMQ_CONFIG_NO_WAIT bool
var EVENT_MAX_RETRY_COUNT int64
var EVENT_REQUEUE_DELAY_MS int64
var FLAG_URL string

var REDIS_HOST string
var REDIS_PORT string
var REDIS_DB int

var SEPARATOR_KEY string
var SERVER_GRACEFUL_SHUTDOWN_TIMEOUT_S int

var TEST bool

var TRACE_OTLP_TARGET string
var TRACE_SAMPLE_RATIO float64

var XMS_API_URI string

var PUBLISH_MAX_RETRY_COUNT int
var PUBLISH_MAX_RETRY_ELAPSED_TIME_SEC int
var PUBLISH_RETRY_WAIT_TIME_SEC int
var PUBLISH_MAX_RETRY_WAIT_TIME_SEC int
var PUBLISH_BACKOFF_FACTOR int

var RABBITMQ_RECONNECTION_DELAY_SECONDS int

// reloadConfig reload config either from file or from system's ENV
// see: infrastructure/configuration/setup.go
func reloadConfig() {
	APP_NAME = viper.GetString("APP_NAME")
	APP_ENV = viper.GetString("APP_ENV")
	APP_PORT = viper.GetInt("APP_PORT")
	APP_URL = viper.GetString("APP_URL")
	FRONTEND_URL = viper.GetString("FRONTEND_URL")
	SERVER_GRACEFUL_SHUTDOWN_TIMEOUT_S = viper.GetInt("SERVER_GRACEFUL_SHUTDOWN_TIMEOUT_S")
	DEBUG = viper.GetBool("DEBUG")
	TEST = viper.GetBool("TEST")
	NOTIFICATION_SERVICE_BASE_URL = viper.GetString("NOTIFICATION_SERVICE_BASE_URL")
	NOTIFICATION_SERVICE_SENDER_EMAIL = viper.GetString("NOTIFICATION_SERVICE_SENDER_EMAIL")
	NOTIFICATION_SERVICE_SENDER_NAME = viper.GetString("NOTIFICATION_SERVICE_SENDER_NAME")
	NR_HOST = viper.GetString("NR_HOST")
	NR_METRIC_PREFIX = viper.GetString("NR_METRIC_PREFIX")
	NR_LICENSE_KEY = viper.GetString("NR_LICENSE_KEY")
	LOG_CALLER = viper.GetBool("LOG_CALLER")
	LOG_LEVEL = viper.GetString("LOG_LEVEL")
	LOG_USE_JSON = viper.GetBool("LOG_USE_JSON")
	DB_POSTGRES_AUTO_MIGRATE = viper.GetBool("DB_POSTGRES_AUTO_MIGRATE")
	DB_MYSQL_AUTO_MIGRATE = viper.GetBool("DB_MYSQL_AUTO_MIGRATE")
	RABBITMQ_CONFIG_DURABLE = viper.GetBool("RABBITMQ_CONFIG_DURABLE")
	RABBITMQ_CONFIG_EXCLUSIVE = viper.GetBool("RABBITMQ_CONFIG_EXCLUSIVE")
	RABBITMQ_CONFIG_AUTO_DELETED = viper.GetBool("RABBITMQ_CONFIG_AUTO_DELETED")
	RABBITMQ_CONFIG_INTERNAL = viper.GetBool("RABBITMQ_CONFIG_INTERNAL")
	RABBITMQ_CONFIG_NO_WAIT = viper.GetBool("RABBITMQ_CONFIG_NO_WAIT")
	EVENT_MAX_RETRY_COUNT = viper.GetInt64("EVENT_MAX_RETRY_COUNT")
	EVENT_REQUEUE_DELAY_MS = viper.GetInt64("EVENT_REQUEUE_DELAY_MS")
	BIFROST_URI = viper.GetString("BIFROST_URI")
	BIFROST_TENANT_CODE = viper.GetString("BIFROST_TENANT_CODE")
	IMAGE_MAX_SIZE_IN_MB = viper.GetInt64("IMAGE_MAX_SIZE_IN_MB")
	IMAGE_MAX_WIDTH = viper.GetInt("IMAGE_MAX_WIDTH")
	IMAGE_MAX_HEIGHT = viper.GetInt("IMAGE_MAX_HEIGHT")
	IMAGE_MIN_WIDTH = viper.GetInt("IMAGE_MIN_WIDTH")
	IMAGE_MIN_HEIGHT = viper.GetInt("IMAGE_MIN_HEIGHT")
	JWT_METHOD = viper.GetString("JWT_METHOD")
	KRATOS_URL = viper.GetString("KRATOS_URL")
	CACHE_WHITELIST_TOKENS = viper.GetString("CACHE_WHITELIST_TOKENS")
	SEPARATOR_KEY = viper.GetString("SEPARATOR_KEY")
	REDIS_DB = viper.GetInt("REDIS_DB")
	REDIS_HOST = viper.GetString("REDIS_HOST")
	XMS_API_URI = viper.GetString("XMS_API_URI")
	FLAG_URL = viper.GetString("FLAG_URL")

	viper.SetDefault("DB_TIMEZONE", "Asia/Bangkok")
	DB_TIMEZONE = viper.GetString("DB_TIMEZONE")

	viper.SetDefault("HTTP_REQUEST_RETRY_COUNT", 3)
	HTTP_REQUEST_RETRY_COUNT = viper.GetInt("HTTP_REQUEST_RETRY_COUNT")

	viper.SetDefault("HTTP_REQUEST_RETRY_WAIT_TIME_MS", 5000)
	HTTP_REQUEST_RETRY_WAIT_TIME_MS = viper.GetInt("HTTP_REQUEST_RETRY_WAIT_TIME_MS")

	viper.SetDefault("HTTP_REQUEST_MAX_WAIT_TIME_MS", 20000)
	HTTP_REQUEST_MAX_WAIT_TIME_MS = viper.GetInt("HTTP_REQUEST_MAX_WAIT_TIME_MS")

	viper.SetDefault("REDIS_PORT", "6379")
	REDIS_PORT = viper.GetString("REDIS_PORT")

	viper.SetDefault("CONTEXT_TIMEOUT_S", 5)
	CONTEXT_TIMEOUT, _ = time.ParseDuration(viper.GetString("CONTEXT_TIMEOUT_S") + "s")

	TRACE_OTLP_TARGET = viper.GetString("TRACE_OTLP_TARGET")
	TRACE_SAMPLE_RATIO = viper.GetFloat64("TRACE_SAMPLE_RATIO")

	PUBLISH_MAX_RETRY_COUNT = viper.GetInt("PUBLISH_MAX_RETRY_COUNT")
	PUBLISH_RETRY_WAIT_TIME_SEC = viper.GetInt("PUBLISH_RETRY_WAIT_TIME_SEC")
	PUBLISH_BACKOFF_FACTOR = viper.GetInt("PUBLISH_BACKOFF_FACTOR")
	PUBLISH_MAX_RETRY_ELAPSED_TIME_SEC = viper.GetInt("PUBLISH_MAX_RETRY_ELAPSED_TIME_SEC")
	PUBLISH_MAX_RETRY_WAIT_TIME_SEC = viper.GetInt("PUBLISH_MAX_RETRY_WAIT_TIME_SEC")

	RABBITMQ_RECONNECTION_DELAY_SECONDS = viper.GetInt("RABBITMQ_RECONNECTION_DELAY_SECONDS")

}

// GetRmqConfig get rabbitmq config from the config vars.
//
// args is optional parameter to specify RabbitMQ config arguments.
func GetRmqConfig(args ...amqp091.Table) rmqConfig.RabbitMQConfig {
	var arg amqp091.Table
	if len(args) > 0 {
		arg = args[0]
	}

	return rmqConfig.RabbitMQConfig{
		AmqpURI:     secret.RABBITMQ_URI,
		Durable:     RABBITMQ_CONFIG_DURABLE,
		Exclusive:   RABBITMQ_CONFIG_EXCLUSIVE,
		AutoDeleted: RABBITMQ_CONFIG_AUTO_DELETED,
		Internal:    RABBITMQ_CONFIG_INTERNAL,
		NoWait:      RABBITMQ_CONFIG_NO_WAIT,
		Arguments:   arg,
	}
}

// GetLogConfig get logging config from the config vars
//
// It used for set up the catalystdk standard log
func GetLogConfig() *log.Config {
	return &log.Config{
		AppName: APP_NAME,
		Caller:  LOG_CALLER,
		Level:   LOG_LEVEL,
		UseJSON: LOG_USE_JSON,
	}
}

// GetTracerConfig get system metric config from the config vars
//
// It used for set up the catalystdk standard tracer for system monitoring
func GetTracerConfig() *traceConfig.Config {
	return &traceConfig.Config{
		AppName:     APP_NAME,
		OtlpTarget:  TRACE_OTLP_TARGET,
		Environment: APP_ENV,
		TraceRatio:  TRACE_SAMPLE_RATIO,
	}
}

// GetFlagConfig get flag config from the config vars
//
// It used for set up the catalystdk flag for feature flagging
func GetFlagConfig() flagConfig.Config {
	return flagConfig.Config{
		Url: FLAG_URL,
	}
}

// GetRedisConfig get redis config config from the config vars
func GetRedisConfig() cache.RedisConfig {
	return cache.RedisConfig{
		Host: REDIS_HOST,
		Port: REDIS_PORT,
		DB:   REDIS_DB,
	}
}
