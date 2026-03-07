package warg

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go.bbkane.com/warg/styles"
	"go.bbkane.com/warg/value"
)

// compactFlagLine holds the pre-computed parts of a flag line for column alignment.
type compactFlagLine struct {
	// leftCol is the flag name/alias/type portion (e.g., "  -e, --editor string")
	leftCol string
	// rightCol is the description + metadata (e.g., "path to editor (default \"vi\") (required) [env: EDITOR]")
	rightCol string
}

// compactBuildFlagLine constructs the left and right columns for a single flag.
func compactBuildFlagLine(s *styles.Styles, name string, f *Flag, val value.Value) compactFlagLine {
	// Build left column: "  -a, --name type" or "      --name type"
	var left strings.Builder
	left.WriteString("  ")
	if f.Alias != "" {
		left.WriteString(s.FlagAlias(f.Alias))
		left.WriteString(", ")
	} else {
		left.WriteString("    ")
	}
	left.WriteString(s.FlagName(name))
	left.WriteString(" ")
	left.WriteString(val.Description())

	// Build right column: description + annotations
	var right strings.Builder
	right.WriteString(f.HelpShort)

	// Add default value
	if val.HasDefault() {
		switch v := val.(type) {
		case value.ScalarValue:
			fmt.Fprintf(&right, " [default: %q]", v.DefaultString())
		case value.SliceValue:
			fmt.Fprintf(&right, " [default: %v]", v.DefaultStringSlice())
		case value.DictValue:
			fmt.Fprintf(&right, " [default: %v]", v.DefaultStringMap())
		}
	}

	// Add required marker
	if f.Required {
		right.WriteString(" [required]")
	}

	// Add env vars
	if len(f.EnvVars) > 0 {
		fmt.Fprintf(&right, " [env: %s]", strings.Join(f.EnvVars, ", "))
	}

	// Add config path
	if f.ConfigPath != "" {
		fmt.Fprintf(&right, " [config: %s]", f.ConfigPath)
	}
	if val.UpdatedBy() != value.UpdatedByUnset {
		fmt.Fprintf(&right, " [setby: %s]", string(val.UpdatedBy()))
		switch v := val.(type) {
		case value.ScalarValue:
			fmt.Fprintf(&right, " [current: %q]", v.String())
		case value.SliceValue:
			fmt.Fprintf(&right, " [current: %v]", v.StringSlice())
		case value.DictValue:
			fmt.Fprintf(&right, " [current: %v]", v.StringMap())
		}
	}

	return compactFlagLine{
		leftCol:  left.String(),
		rightCol: right.String(),
	}
}

// compactVisibleLen returns the visible length of a string, stripping ANSI escape sequences.
func compactVisibleLen(s string) int {
	n := 0
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		n++
	}
	return n
}

// compactWrapText wraps text at word boundaries to fit within maxWidth characters.
// Each continuation line is indented by indent spaces.
func compactWrapText(text string, indent int, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	indentStr := strings.Repeat(" ", indent)
	var result strings.Builder
	lineLen := 0

	for i, word := range words {
		wordLen := len(word)
		if i == 0 {
			result.WriteString(word)
			lineLen = wordLen
			continue
		}

		// Check if adding this word would exceed the max width.
		if lineLen+1+wordLen > maxWidth {
			result.WriteString("\n")
			result.WriteString(indentStr)
			result.WriteString(word)
			lineLen = indent + wordLen
		} else {
			result.WriteString(" ")
			result.WriteString(word)
			lineLen += 1 + wordLen
		}
	}

	return result.String()
}

// compactPrintFlags prints a set of flag lines with aligned columns, respecting terminal width.
func compactPrintFlags(p *styles.Printer, lines []compactFlagLine, termWidth int) {
	if len(lines) == 0 {
		return
	}

	// Find the maximum left column width for alignment
	maxLeftWidth := 0
	for _, line := range lines {
		w := compactVisibleLen(line.leftCol)
		if w > maxLeftWidth {
			maxLeftWidth = w
		}
	}

	// Add a gutter of at least 3 spaces between columns
	const gutter = 3
	descCol := maxLeftWidth + gutter

	for _, line := range lines {
		leftWidth := compactVisibleLen(line.leftCol)
		padding := strings.Repeat(" ", descCol-leftWidth)

		if termWidth > 0 {
			// Available width for description text
			availWidth := termWidth - descCol
			if availWidth < 20 {
				// If the terminal is very narrow, just print without wrapping
				p.Printf("%s%s%s\n", line.leftCol, padding, line.rightCol)
			} else {
				wrapped := compactWrapText(line.rightCol, descCol, availWidth)
				p.Printf("%s%s%s\n", line.leftCol, padding, wrapped)
			}
		} else {
			p.Printf("%s%s%s\n", line.leftCol, padding, line.rightCol)
		}
	}
}

// compactCmdHelp returns an Action that prints Compact-style help for the current command.
func compactCmdHelp() Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		s, err := conditionallyEnableStyle(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		p := styles.NewPrinter(f)
		cur := cmdCtx.ParseState.CurrentCmd
		termWidth := TermWidth(cmdCtx.Stdout, cmdCtx.Flags)

		// Build the usage path: <app> [section...] <cmd>
		var usagePath strings.Builder
		usagePath.WriteString(string(cmdCtx.App.Name))
		for _, sec := range cmdCtx.ParseState.SectionPath {
			usagePath.WriteString(" ")
			usagePath.WriteString(sec)
		}
		usagePath.WriteString(" ")
		usagePath.WriteString(cmdCtx.ParseState.CurrentCmdName)

		// Usage line
		p.Printf("%s:\n", s.Header("Usage"))
		if cur.AllowForwardedArgs {
			p.Printf("  %s [flags] -- [args]\n", usagePath.String())
		} else {
			p.Printf("  %s [flags]\n", usagePath.String())
		}
		p.Println()

		// Description
		if cur.HelpLong != "" {
			p.Println(cur.HelpLong)
		} else {
			p.Println(cur.HelpShort)
		}
		p.Println()

		// Command Flags
		var cmdFlagLines []compactFlagLine
		for _, name := range cmdCtx.ParseState.CurrentCmd.Flags.SortedNames() {
			fl := cmdCtx.ParseState.CurrentCmd.Flags[name]
			val := cmdCtx.ParseState.FlagValues[name]
			cmdFlagLines = append(cmdFlagLines, compactBuildFlagLine(&s, name, &fl, val))
		}
		if len(cmdFlagLines) > 0 {
			p.Printf("%s:\n", s.Header("Flags"))
			compactPrintFlags(p, cmdFlagLines, termWidth)
			p.Println()
		}

		// Global Flags
		var globalFlagLines []compactFlagLine
		for _, name := range cmdCtx.App.GlobalFlags.SortedNames() {
			fl := cmdCtx.App.GlobalFlags[name]
			val := cmdCtx.ParseState.FlagValues[name]
			globalFlagLines = append(globalFlagLines, compactBuildFlagLine(&s, name, &fl, val))
		}
		if len(globalFlagLines) > 0 {
			p.Printf("%s:\n", s.Header("Global Flags"))
			compactPrintFlags(p, globalFlagLines, termWidth)
			p.Println()
		}

		// Footer
		if cur.Footer != "" {
			p.Printf("%s:\n", s.Header("Footer"))
			p.Println(cur.Footer)
		}

		return nil
	}
}

// compactSectionHelp returns an Action that prints Compact-style help for the current section.
func compactSectionHelp() Action {
	return func(cmdCtx CmdContext) error {
		file := cmdCtx.Stdout
		f := bufio.NewWriter(file)
		defer f.Flush()

		s, err := conditionallyEnableStyle(cmdCtx.Flags, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error enabling color. Continuing without: %v\n", err)
		}

		p := styles.NewPrinter(f)
		cur := cmdCtx.ParseState.CurrentSection
		termWidth := TermWidth(cmdCtx.Stdout, cmdCtx.Flags)

		// Build the usage path: <app> [section...]
		var usagePath strings.Builder
		usagePath.WriteString(string(cmdCtx.App.Name))
		for _, sec := range cmdCtx.ParseState.SectionPath {
			usagePath.WriteString(" ")
			usagePath.WriteString(sec)
		}

		// Usage line
		p.Printf("%s:\n", s.Header("Usage"))
		p.Printf("  %s [command]\n", usagePath.String())
		p.Println()

		// Description
		if cur.HelpLong != "" {
			p.Println(cur.HelpLong)
		} else {
			p.Println(cur.HelpShort)
		}
		p.Println()

		// Available Commands
		if len(cur.Cmds) > 0 {
			p.Printf("%s:\n", s.Header("Available Commands"))

			// Compute max command name length for alignment
			maxNameLen := 0
			for _, k := range cur.Cmds.SortedNames() {
				if len(k) > maxNameLen {
					maxNameLen = len(k)
				}
			}

			const gutter = 3
			descCol := 2 + maxNameLen + gutter // 2 for indent

			for _, k := range cur.Cmds.SortedNames() {
				name := s.CommandName(k)
				nameVisLen := compactVisibleLen(name)
				padding := strings.Repeat(" ", descCol-2-nameVisLen+gutter-gutter) // align to descCol from after the indent
				// Recalculate: indent(2) + name + padding + desc
				pad := strings.Repeat(" ", 2+maxNameLen+gutter-2-nameVisLen)

				desc := cur.Cmds[string(k)].HelpShort
				if termWidth > 0 {
					availWidth := termWidth - descCol
					if availWidth >= 20 {
						desc = compactWrapText(desc, descCol, availWidth)
					}
				}
				_ = padding
				p.Printf("  %s%s%s\n", name, pad, desc)
			}
			p.Println()
		}

		// Sections (sub-sections)
		if len(cur.Sections) > 0 {
			p.Printf("%s:\n", s.Header("Additional Commands"))

			maxNameLen := 0
			for _, k := range cur.Sections.SortedNames() {
				if len(k) > maxNameLen {
					maxNameLen = len(k)
				}
			}

			const gutter = 3
			descCol := 2 + maxNameLen + gutter

			for _, k := range cur.Sections.SortedNames() {
				name := s.SectionName(k)
				nameVisLen := compactVisibleLen(name)
				pad := strings.Repeat(" ", 2+maxNameLen+gutter-2-nameVisLen)

				desc := cur.Sections[k].HelpShort
				if termWidth > 0 {
					availWidth := termWidth - descCol
					if availWidth >= 20 {
						desc = compactWrapText(desc, descCol, availWidth)
					}
				}
				p.Printf("  %s%s%s\n", name, pad, desc)
			}
			p.Println()
		}

		// Global Flags
		var globalFlagLines []compactFlagLine
		for _, name := range cmdCtx.App.GlobalFlags.SortedNames() {
			fl := cmdCtx.App.GlobalFlags[name]
			val := cmdCtx.ParseState.FlagValues[name]
			globalFlagLines = append(globalFlagLines, compactBuildFlagLine(&s, name, &fl, val))
		}
		if len(globalFlagLines) > 0 {
			p.Printf("%s:\n", s.Header("Global Flags"))
			compactPrintFlags(p, globalFlagLines, termWidth)
			p.Println()
		}

		// Footer hint
		p.Printf("Use \"%s [command] --help\" for more information about a command.\n", usagePath.String())

		// Section footer
		if cur.Footer != "" {
			p.Println()
			p.Printf("%s:\n", s.Header("Footer"))
			p.Println(cur.Footer)
		}

		return nil
	}
}
