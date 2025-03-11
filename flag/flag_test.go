package flag_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/cli"
	"go.bbkane.com/warg/flag"
)

func TestFlagMap_SortedNames(t *testing.T) {
	emptyFlag := flag.NewFlag("", nil)

	fm := cli.FlagMap{
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
