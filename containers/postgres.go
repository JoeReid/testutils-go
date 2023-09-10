package containers

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

func Postgres(t *testing.T, opts ...PostgresOpt) *sqlx.DB {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("failed to construct docker pool: %v", err)
	}

	if err := pool.Client.Ping(); err != nil {
		t.Fatalf("failed to connect to docker: %v", err)
	}

	postgresOpts := &postgresOpts{
		tag:           "latest",
		migrationsDir: "",
	}
	for _, opt := range opts {
		opt(postgresOpts)
	}

	resource, err := postgresOpts.run(pool)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	db, err := postgresOpts.sqlx(pool, resource)
	if err != nil {
		t.Fatalf("failed to connect to postgres container: %v", err)
	}

	if err := postgresOpts.migrate(pool, resource, db.DB); err != nil {
		t.Fatalf("failed to migrate postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("failed to purge postgres container: %v", err)
		}
	})

	return db
}

type PostgresOpt func(*postgresOpts)

func WithPostgresTag(tag string) PostgresOpt {
	return func(o *postgresOpts) {
		o.tag = tag
	}
}

func WithMigrationDir(dir string) PostgresOpt {
	return func(o *postgresOpts) {
		o.migrationsDir = dir
	}
}

type postgresOpts struct {
	tag           string
	migrationsDir string
}

func (o *postgresOpts) run(pool *dockertest.Pool) (*dockertest.Resource, error) {
	return pool.Run("postgres", o.tag, []string{
		"POSTGRES_DB=postgres",
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=postgres",
	})
}

func (o *postgresOpts) sqlx(pool *dockertest.Pool, resource *dockertest.Resource) (*sqlx.DB, error) {
	var db *sqlx.DB

	connect := func() (err error) {
		data := strings.Join([]string{
			"host=localhost",
			"port=" + resource.GetPort("5432/tcp"),
			"user=postgres",
			"password=postgres",
			"dbname=postgres",
			"sslmode=disable",
		}, " ")

		db, err = sqlx.Open("postgres", data)
		if err != nil {
			return err
		}
		return db.Ping()
	}

	return db, pool.Retry(connect)
}

func (o *postgresOpts) migrate(pool *dockertest.Pool, resource *dockertest.Resource, db *sql.DB) error {
	if o.migrationsDir == "" {
		return nil
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+o.migrationsDir, "postgres", driver)
	if err != nil {
		return err
	}

	return m.Up()
}
