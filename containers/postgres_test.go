package containers

import (
	"context"
	"testing"
	"time"
)

func TestPostgres(t *testing.T) {
	t.Parallel()

	db := Postgres(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping postgres: %v", err)
	}
}

func TestPostgres_Migrations(t *testing.T) {
	t.Parallel()

	db := Postgres(t, WithPostgresMigrations("testdata/postgres_migrations"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	if _, err := db.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test')"); err != nil {
		t.Fatalf("failed to insert into postgres: %v", err)
	}
}
