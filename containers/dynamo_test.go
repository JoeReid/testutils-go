package containers

import (
	"context"
	"testing"
	"time"
)

func TestDynamoDB(t *testing.T) {
	t.Parallel()

	db := DynamoDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	_, err := db.ListTables().AllWithContext(ctx)
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}
}
