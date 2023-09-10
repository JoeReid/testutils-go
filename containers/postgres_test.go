package containers

import (
	"context"
	"testing"
	"time"
)

func TestPostgres(t *testing.T) {
	db := Postgres(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping postgres: %v", err)
	}
}
