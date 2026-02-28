package warg

import (
	"os"
	"sort"

	"github.com/mattn/go-isatty"
	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/styles"
	"go.bbkane.com/warg/value/scalar"
)

func DefaultHelpCmdMap() CmdMap {
	allCmdsHelp := NewCmd("", buildHelpAction(detailedCmdHelp(), allCmdsSectionHelp()))
	return CmdMap{
		"default":     allCmdsHelp,
		"detailed":    NewCmd("", buildHelpAction(detailedCmdHelp(), detailedSectionHelp())),
		"outline":     outlineHelp(),
		"allcommands": allCmdsHelp,
	}
}

func DefaultHelpFlagMap(defaultChoice string, choices []string) FlagMap {
	return FlagMap{
		"--help": NewFlag(
			"Print help",
			scalar.String(
				scalar.Choices(choices...),
				scalar.Default(defaultChoice),
			),
			Alias("-h"),
		),
	}
}

func buildHelpAction(cmdAction Action, secAction Action) Action {
	return func(cmdCtx CmdContext) error {
		com := cmdCtx.ParseState.CurrentCmd
		if com != nil {
			return cmdAction(cmdCtx)
		}
		return secAction(cmdCtx)
	}
}

// leftPad pads a string `s` with pad `pad` `plength` times
//
// In Python: (pad * plength) + s
func leftPad(s string, pad string, plength int) string {
	// https://stackoverflow.com/a/45456649/2958070
	for i := 0; i < plength; i++ {
		s = pad + s
	}
	return s
}

// ColorEnabled looks for a passed --color flag with an underlying string value. If
// ConditionallyEnableColor2 looks for a passed --color flag with an underlying string value. If
// it exists and is set to "true", or if it exists, is set to "auto",
// and the passed file is a TTY, an enabled Color is returned.
func ColorEnabled(pf PassedFlags, file *os.File) bool {
	// default to trying to use color
	useColor := "auto"
	// respect a --color string
	if useColorI, exists := pf["--color"]; exists {
		if useColorUnder, isStr := useColorI.(string); isStr {
			useColor = useColorUnder
		}
	}

	startEnabled := useColor == "true" || (useColor == "auto" && isatty.IsTerminal(file.Fd()))
	return startEnabled
}

// conditionallyEnableStyle looks for a passed --color flag with an underlying string value. If
// it exists and is set to "true", or if it exists, is set to "auto",
// and the passed file is a TTY, an enabled Styles is returned.
func conditionallyEnableStyle(pf PassedFlags, file *os.File) (styles.Styles, error) {
	startEnabled := ColorEnabled(pf, file)
	if !startEnabled {
		return styles.NewEmptyStyles(), nil
	}
	err := gocolor.EnableConsole()
	if err != nil {
		return styles.NewEmptyStyles(), err
	}
	return styles.NewEnabledStyles(), nil
}

// sortedKeys returns the keys of the map m in sorted order.
// copied and modified from https://cs.opensource.google/go/x/exp/+/master:maps/maps.go;l=10;drc=79cabaa25d7518588d46eb676385c8dff49670c3
func sortedKeys[M ~map[string]V, V any](m M) []string {
	r := make([]string, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}
