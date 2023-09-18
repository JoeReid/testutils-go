package containers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	t.Parallel()

	client := NewRedis(t).GoRedisClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := client.Set(ctx, "test", "test", 0).Err()
	require.NoError(t, err)

	val, err := client.Get(ctx, "test").Result()
	require.NoError(t, err)
	assert.Equal(t, "test", val)
}
