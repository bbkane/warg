package flag_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/wargcore"
)

func TestFlagMap_SortedNames(t *testing.T) {
	emptyFlag := flag.NewFlag("", nil)

	fm := wargcore.FlagMap{
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
