package server

import (
	"context"
	"time"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/cache"

	"github.com/etherlabsio/healthcheck/v2"
	"github.com/go-chi/chi/v5"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/database"
)

// addHealthCheck added healthcheck http routes to validate the application stack healthiness.
//
// It check connection of postgres databases, redis, and mongodb
func addHealthCheck(r *chi.Mux, dbConnections database.DatabaseConnections, cache cache.Cache) {
	r.Get("/health", healthcheck.HandlerFunc(
		healthcheck.WithTimeout(5*time.Second),
		healthcheck.WithChecker(
			"postgres-master", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					db, err := dbConnections.Postgres[database.SamplePostgresDBName].Master.DB()
					if err != nil {
						return err
					}

					return db.PingContext(ctx)
				},
			),
		),
		healthcheck.WithChecker(
			"postgres-slave", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					db, err := dbConnections.Postgres[database.SamplePostgresDBName].Slave.DB()
					if err != nil {
						return err
					}

					return db.PingContext(ctx)
				},
			),
		),
		healthcheck.WithChecker(
			"mongo", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return dbConnections.Mongo[database.SampleMongoDBName].DB.Client().Ping(ctx, nil)
				},
			),
		),
		healthcheck.WithChecker(
			"redis", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return cache.PING(ctx)
				},
			),
		),
	))
}
