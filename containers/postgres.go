package containers

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

type Postgres struct {
	resource *dockertest.Resource
}

func (p *Postgres) SQLX(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open("postgres", strings.Join([]string{
		"host=localhost",
		"port=" + p.resource.GetPort("5432/tcp"),
		"user=postgres",
		"password=postgres",
		"dbname=postgres",
		"sslmode=disable",
	}, " "))
	if err != nil {
		t.Fatal("failed to connect to postgres", err)
	}

	return db
}

func (p *Postgres) DB(t *testing.T) *sql.DB {
	t.Helper()

	return p.SQLX(t).DB
}

func (p *Postgres) Migrate(t *testing.T, fileDir string) *migrate.Migrate {
	t.Helper()

	driver, err := postgres.WithInstance(p.DB(t), &postgres.Config{})
	if err != nil {
		t.Fatal("failed to create migration driver", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+fileDir, "postgres", driver)
	if err != nil {
		t.Fatal("failed to create migration instance", err)
	}

	return m
}

func NewPostgres(t *testing.T) *Postgres {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to construct docker pool: %v", err)
	}

	if err := pool.Client.Ping(); err != nil {
		t.Fatalf("failed to connect to docker: %v", err)
	}

	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_DB=postgres",
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=postgres",
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	r := &Postgres{resource: resource}

	if err := pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		return r.SQLX(t).PingContext(ctx)
	}); err != nil {
		t.Fatalf("failed to connect to postgres container: %v", err)
	}

	return r
}
