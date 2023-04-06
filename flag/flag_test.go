package flag_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/flag"
)

func TestFlagMap_SortedNames(t *testing.T) {
	emptyFlag := flag.New("", nil)

	fm := flag.FlagMap{
		"c": emptyFlag,
		"a": emptyFlag,
		"d": emptyFlag,
		"b": emptyFlag,
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
