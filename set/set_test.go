package set_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/set"
)

func TestSet(t *testing.T) {
	require := require.New(t)
	s := set.New[string]()

	require.False(s.Contains("a"))
	s.Add("a")
	require.True(s.Contains("a"))
	s.Delete("a")
	require.False(s.Contains("a"))
}
