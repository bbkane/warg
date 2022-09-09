package columns

import (
	"fmt"
	"os"
	"strings"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/help/common"
)

const flagHelpSep = " : "
const flagIndent = "  "
const flagAliasNameSep = ", "

// padding returns the spaces to pad a name to fit a given length
func padding(n flag.Name, length int) string {
	lenFlagName := len(n)
	if lenFlagName > length {
		panic("lenFlagName > length")
	}
	return strings.Repeat(" ", length-lenFlagName)
}

// maxFlagColWidth calculates the max width of everything in the flag column.
// NOTE: flagIndent is just considered it's own column
func maxFlagColWidth(fm flag.FlagMap) int {
	m := 0
	for name, fl := range fm {
		lenNameAliasOptionalSep := len(name)
		if fl.Alias != "" {
			lenNameAliasOptionalSep = lenNameAliasOptionalSep + len(fl.Alias) + len(flagAliasNameSep)
		}
		if lenNameAliasOptionalSep > m {
			m = lenNameAliasOptionalSep
		}
	}
	return m
}

func printNoSpace(s ...any) {
	for _, si := range s {
		fmt.Print(si)
	}
}

func printlnNoSpace(s ...any) {
	printNoSpace(s...)
	fmt.Println()
}

func NoWrapCommandHelp(file *os.File, cur *command.Command, hi common.HelpInfo) command.Action {
	return func(ctx command.Context) error {
		maxFlagColWidth_ := maxFlagColWidth(hi.AvailableFlags)
		// fmt.Println(flagIndent + strings.Repeat("_", maxFlagColWidth_))
		for _, flagName := range hi.AvailableFlags.SortedNames() {
			fl := hi.AvailableFlags[flagName]
			fmt.Print(flagIndent)
			paddingWidth := maxFlagColWidth_
			// Adjust padding if we get an alias
			if string(fl.Alias) != "" {
				printNoSpace(
					string(fl.Alias),
					flagAliasNameSep,
				)
				paddingWidth = paddingWidth - len(fl.Alias) - len(flagAliasNameSep)
			}
			printlnNoSpace(
				string(flagName),
				padding(flagName, paddingWidth),
				flagHelpSep,
				string(fl.HelpShort),
			)
		}
		return nil
	}
}
