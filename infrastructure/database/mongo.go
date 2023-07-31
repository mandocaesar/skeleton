package database

import (
	"context"
	"fmt"
	"net/url"

	"github.com/machtwatch/catalystdk/go/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfigs represent collection of Mongo database configuration
type MongoConfigs map[string]MongoConfig

// MongoConnections represent collection of Mongo database connections
type MongoConnections map[string]MongoConnection

// MongoConfig represent Mongo database configuration
type MongoConfig struct {
	Name     string
	Timeout  int
	DBname   string
	Username string
	Password string
	Host     string
	Port     string
}

// MongoConnections represent Mongo database connections
type MongoConnection struct {
	Name string
	DB   *mongo.Database
}

func openMongo(ctx context.Context, config MongoConfig) *mongo.Database {

	var userInfo *url.Userinfo
	if config.Username != "" {
		userInfo = url.UserPassword(config.Username, config.Password)
	}

	URL := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		User:   userInfo,
	}

	mongoClientOpts := options.Client().ApplyURI(URL.String())
	mongoClientOpts.SetDirect(true)
	mongoClientOpts.SetRetryWrites(false)

	client, err := mongo.Connect(ctx, mongoClientOpts)
	if err != nil {
		log.StdFatal(ctx, map[string]interface{}{"host": config.Host, "port": config.Port}, err, "mongo.Connect() monggodb connection got error - openMongo()")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.StdFatal(ctx, map[string]interface{}{"host": config.Host, "port": config.Port}, err, "client.Ping() mongodb client ping got error - openMongo()")
	}

	log.Info("successfully connected to mongodb")

	return client.Database(config.DBname)
}

func closeMongo(ctx context.Context, db *mongo.Database) {
	if err := db.Client().Disconnect(context.Background()); err != nil {
		log.StdFatal(ctx, nil, err, "closeMongo() closing mongodb database got error")
	}
}
