package flag_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/flag"
)

func TestFlagMap_SortedNames(t *testing.T) {
	emptyFlag := flag.Flag{
		Alias:                 "",
		ConfigPath:            "",
		EnvVars:               nil,
		EmptyValueConstructor: nil,
		HelpShort:             "",
		Required:              false,
		IsCommandFlag:         false,
		SetBy:                 "",
		Value:                 nil,
	}
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
