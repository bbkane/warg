package metadata_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/metadata"
)

type testKey struct{}

func TestMetadata(t *testing.T) {
	require := require.New(t)

	md := metadata.New("key1", "value1", testKey{}, 12345)

	value, exists := md.Get("key1")
	require.True(exists)
	require.Equal("value1", value)

	value = md.MustGet(testKey{})
	require.Equal(12345, value)

	require.Panics(func() {
		md.MustGet("not there")
	})
}

func TestMetadata_Empty(t *testing.T) {
	require := require.New(t)

	md := metadata.Empty()

	_, exists := md.Get("anykey")
	require.False(exists)
}
