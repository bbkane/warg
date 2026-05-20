package warg

import (
	"os"
	"sort"
	"strconv"

	"github.com/mattn/go-isatty"
	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/styles"
	"go.bbkane.com/warg/value/scalar"
	"golang.org/x/term"
)

// DefaultHelpCmdMap returns the built-in help command implementations: "default", "detailed",
// "outline", "allcommands", and "compact".
func DefaultHelpCmdMap() CmdMap {
	allCmdsHelp := NewCmd("", buildHelpAction(detailedCmdHelp(), allCmdsSectionHelp()))
	return CmdMap{
		"default":     NewCmd("", buildHelpAction(compactCmdHelp(), allCmdsSectionHelp())),
		"detailed":    NewCmd("", buildHelpAction(detailedCmdHelp(), detailedSectionHelp())),
		"outline":     outlineHelp(),
		"allcommands": allCmdsHelp,
		"compact":     NewCmd("", buildHelpAction(compactCmdHelp(), compactSectionHelp())),
	}
}

// DefaultHelpFlagMap returns a [FlagMap] containing a "--help" / "-h" flag
// with the given default choice and available choices.
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

// ColorEnabled reports whether ANSI color output should be enabled based on the
// --color flag value in pf and whether file is a terminal.
// Returns true if --color is "true", or if --color is "auto" and file is a TTY.
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

// TermWidth determines the terminal width to use for wrapping help output.
// Returns 0 for infinite width. Respects the --term-width flag if present:
// "auto" detects terminal size (falls back to 0), "infinite" returns 0,
// and a positive integer is used directly.
func TermWidth(file *os.File, pf PassedFlags) int {
	termWidthI, exists := pf["--term-width"]
	if !exists {
		return 0
	}
	switch termWidthI.(string) {
	case "infinite":
		return 0
	case "auto":
		fd := int(file.Fd())
		if !term.IsTerminal(fd) {
			return 0
		}
		termWidth, _, err := term.GetSize(fd)
		if err != nil {
			return 0
		}
		return termWidth
	default:
		// try to parse as positive integer
		var err error
		termWidth, err := strconv.Atoi(termWidthI.(string))
		if err != nil || termWidth < 1 {
			return 0
		}
		return termWidth
	}
}

// conditionallyEnableStyle looks for a passed --color flag with an underlying string value. If
// it exists and is set to "true", or if it exists, is set to "auto",
// and the passed file is a TTY, an enabled Styles is returned. Disable with the disable param, so we can use this before parsing flags.
func conditionallyEnableStyle(disable bool, pf PassedFlags, file *os.File) (styles.Styles, error) {
	if disable {
		return styles.NewEmptyStyles(), nil
	}
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
