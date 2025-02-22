package testing

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"

	"github.com/algorand/indexer/idb/postgres/internal/schema"
)

var testpg = flag.String(
	"test-pg", "", "postgres connection string; resets the database")

// SetupPostgres starts a gnomock postgres DB then returns the database object,
// the connection string and a shutdown function.
func SetupPostgres(t *testing.T) (*pgxpool.Pool, string, func()) {
	if testpg != nil && *testpg != "" {
		// use non-docker Postgresql
		connStr := *testpg

		db, err := pgxpool.Connect(context.Background(), connStr)
		require.NoError(t, err, "Error opening postgres connection")

		_, err = db.Exec(
			context.Background(), `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
		require.NoError(t, err)

		shutdownFunc := func() {
			db.Close()
		}

		return db, connStr, shutdownFunc
	}

	p := postgres.Preset(
		postgres.WithVersion("12.5"),
		postgres.WithUser("gnomock", "gnomick"),
		postgres.WithDatabase("mydb"),
	)
	container, err := gnomock.Start(p)
	require.NoError(t, err, "Error starting gnomock")

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable",
		container.Host, container.DefaultPort(),
		"gnomock", "gnomick", "mydb",
	)

	db, err := pgxpool.Connect(context.Background(), connStr)
	require.NoError(t, err, "Error opening postgres connection")

	shutdownFunc := func() {
		db.Close()
		err = gnomock.Stop(container)
		require.NoError(t, err, "Error stoping gnomock")
	}

	return db, connStr, shutdownFunc
}

// SetupPostgresWithSchema is equivalent to SetupPostgres() but also creates the
// indexer schema.
func SetupPostgresWithSchema(t *testing.T) (*pgxpool.Pool, string, func()) {
	db, connStr, shutdownFunc := SetupPostgres(t)

	_, err := db.Exec(context.Background(), schema.SetupPostgresSql)
	require.NoError(t, err)

	return db, connStr, shutdownFunc
}
