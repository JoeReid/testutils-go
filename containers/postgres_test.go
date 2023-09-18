package containers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {
	t.Parallel()

	psql := NewPostgres(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	err := psql.SQLX(t).PingContext(ctx)
	require.NoError(t, err)
}

func TestPostgres_Migrations(t *testing.T) {
	t.Parallel()

	psql := NewPostgres(t)

	err := psql.Migrate(t, "testdata/postgres_migrations").Up()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	_, err = psql.SQLX(t).ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test')")
	assert.NoError(t, err)
}
