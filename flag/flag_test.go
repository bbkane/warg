package flag_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/flag"
)

func TestFlagMap_SortedNames(t *testing.T) {
	fm := flag.FlagMap{
		"c": flag.Flag{},
		"a": flag.Flag{},
		"d": flag.Flag{},
		"b": flag.Flag{},
	}
	require.Equal(
		t,
		[]flag.Name{
			flag.Name("a"),
			flag.Name("b"),
			flag.Name("c"),
			flag.Name("d"),
		},
		fm.SortedNames(),
	)
}
