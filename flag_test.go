package warg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg"
)

func TestFlagMap_SortedNames(t *testing.T) {
	emptyFlag := warg.NewFlag("", nil)

	fm := warg.FlagMap{
		"c": emptyFlag,
		"a": emptyFlag,
		"d": emptyFlag,
		"b": emptyFlag,
	}
	require.Equal(
		t,
		[]string{
			string("a"),
			string("b"),
			string("c"),
			string("d"),
		},
		fm.SortedNames(),
	)
}
