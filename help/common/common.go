package common

import (
	"os"
	"sort"

	"github.com/mattn/go-isatty"
	"go.bbkane.com/gocolor"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
)

// HelpInfo lists common information available to a help function
type HelpInfo struct {

	// AvailableFlags for the current section or commmand, including inherted flags from parent sections.
	// All flags are Resolved if possible (i.e., flag.SetBy != "")
	AvailableFlags flag.FlagMap
	// RootSection of the app. Especially useful for printing all sections and commands
	RootSection section.SectionT
}

// LeftPad pads a string `s` with pad `pad` `plength` times
//
// In Python: (pad * plength) + s
func LeftPad(s string, pad string, plength int) string {
	// https://stackoverflow.com/a/45456649/2958070
	for i := 0; i < plength; i++ {
		s = pad + s
	}
	return s
}

// ConditionallyEnableColor looks for a passed --color flag with an underlying string value. If
// it exists and is set to "true", or if it exists, is set to "auto",
// and the passed file is a TTY, an enabled Color is returned.
func ConditionallyEnableColor(pf command.PassedFlags, file *os.File) (gocolor.Color, error) {
	// default to trying to use color
	useColor := "auto"
	// respect a --color string
	if useColorI, exists := pf["--color"]; exists {
		if useColorUnder, isStr := useColorI.(string); isStr {
			useColor = useColorUnder
		}
	}

	startEnabled := useColor == "true" || (useColor == "auto" && isatty.IsTerminal(file.Fd()))
	return gocolor.Prepare(startEnabled)

}

func FmtHeader(col *gocolor.Color, header string) string {
	return col.Add(col.Bold+col.Underline, header)
}

func FmtSectionName(col *gocolor.Color, sectionName section.Name) string {
	return col.Add(col.Bold+col.FgCyan, string(sectionName))
}

func FmtCommandName(col *gocolor.Color, commandName command.Name) string {
	return col.Add(col.Bold+col.FgGreen, string(commandName))
}

func FmtFlagName(col *gocolor.Color, flagName flag.Name) string {
	return col.Add(col.Bold+col.FgYellow, string(flagName))
}

func FmtFlagAlias(col *gocolor.Color, flagAlias flag.Name) string {
	return col.Add(col.Bold+col.FgYellow, string(flagAlias))
}

// SortedKeys returns the keys of the map m in sorted order.
// copied and modified from https://cs.opensource.google/go/x/exp/+/master:maps/maps.go;l=10;drc=79cabaa25d7518588d46eb676385c8dff49670c3
func SortedKeys[M ~map[string]V, V any](m M) []string {
	r := make([]string, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}
