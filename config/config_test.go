package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("no env variables", func(t *testing.T) {
		os.Clearenv()
		config := New()

		require.Equal(t, serverAddressDefault, config.Server.Address)
		require.Equal(t, serverWriteTimeoutDefault, config.Server.WriteTimeout)
		require.Equal(t, serverReadTimeoutDefault, config.Server.ReadTimeout)
		require.Equal(t, serverIdleTimeoutDefault, config.Server.IdleTimeout)
		require.Equal(t, cacheCapacityDefault, config.Cache.Capacity)
	})

	t.Run("with env variables", func(t *testing.T) {
		os.Setenv("SERVER_ADDRESS", "new-address.com:8080")
		os.Setenv("SERVER_WRITE_TIMEOUT", "100")
		os.Setenv("SERVER_READ_TIMEOUT", "500")
		os.Setenv("SERVER_IDLE_TIMEOUT", "5")
		os.Setenv("CACHE_CAPACITY", "10000000000000000")

		config := New()

		require.Equal(t, "new-address.com:8080", config.Server.Address)
		require.Equal(t, time.Second*100, config.Server.WriteTimeout)
		require.Equal(t, time.Second*500, config.Server.ReadTimeout)
		require.Equal(t, time.Second*5, config.Server.IdleTimeout)
		require.Equal(t, 10000000000000000, config.Cache.Capacity)
	})
}
