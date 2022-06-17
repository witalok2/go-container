package postgres

import (
	"context"
	"testing"
)

func TestPostgres(t *testing.T) {
	pgsql := NewPostgresContainer()
	pgsql.RunMigrations("../test/scripts/migrations")
	pgsql.ExecuteSQLFile("../test/scripts/testdata/testdata.sql")

	var uuid int
	err := pgsql.postgresDB.QueryRow(context.Background(), `SELECT COUNT(id) FROM users`).Scan(&uuid)
	if err != nil {
		ErrPurge(t, pgsql)
	}

	if uuid == 0 {
		ErrPurge(t, pgsql)
	}

	pgsql.Purge()
}

func ErrPurge(t *testing.T, container *PostgresContainer) {
	container.Purge()
	t.Error()
}
