package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewDefaultCache(t *testing.T) {
	var dc = NewDefaultCache(time.Minute)

	t.Run("implements ValidationKeysCache Interface", func(t *testing.T) {
		assert.Implements(t, (*ValidationKeysCache)(nil), dc)
	})

	t.Run("shouldn't find any key in empty cache", func(t *testing.T) {
		_, found := dc.Get("test")
		require.False(t, found)
	})

	t.Run("should get cached value", func(t *testing.T) {
		dc.Set("test", "42")
		cached, found := dc.Get("test")
		require.True(t, found)
		require.Equal(t, cached, "42")
	})
}
