// Package styles defines ANSI color styles used by warg's help output and error formatting.
package styles

import (
	"fmt"
	"io"

	"go.bbkane.com/gocolor"
)

// Styles holds ANSI color codes for different UI elements in help and error output.
// Use [NewEnabledStyles] for colored output or [NewEmptyStyles] to disable colors.
type Styles struct {
	CommandNameCode gocolor.Code
	ErrorAltCode    gocolor.Code
	ErrorCode       gocolor.Code
	FlagAliasCode   gocolor.Code
	FlagNameCode    gocolor.Code
	HeaderCode      gocolor.Code
	LabelCode       gocolor.Code
	SectionNameCode gocolor.Code

	// DefaultCode should probably always be to gocolor.Default or gocolor.Empty
	DefaultCode gocolor.Code
}

// NewEmptyStyles returns a [Styles] with no ANSI codes (no color output).
func NewEmptyStyles() Styles {
	return Styles{
		CommandNameCode: gocolor.Empty,
		ErrorAltCode:    gocolor.Empty,
		ErrorCode:       gocolor.Empty,
		FlagAliasCode:   gocolor.Empty,
		FlagNameCode:    gocolor.Empty,
		HeaderCode:      gocolor.Empty,
		LabelCode:       gocolor.Empty,
		SectionNameCode: gocolor.Empty,
		DefaultCode:     gocolor.Empty,
	}
}

// NewEnabledStyles returns a [Styles] with default ANSI color codes for terminal output.
func NewEnabledStyles() Styles {
	return Styles{
		CommandNameCode: gocolor.Bold + gocolor.FgGreen,
		ErrorAltCode:    gocolor.Bold + gocolor.FgYellowBright,
		ErrorCode:       gocolor.Bold + gocolor.FgRedBright,
		FlagAliasCode:   gocolor.Bold + gocolor.FgYellow,
		FlagNameCode:    gocolor.Bold + gocolor.FgYellow,
		HeaderCode:      gocolor.Bold + gocolor.Underline,
		LabelCode:       gocolor.Bold,
		SectionNameCode: gocolor.Bold + gocolor.FgCyan,
		DefaultCode:     gocolor.Default,
	}
}

// CommandName wraps v with the command name style.
func (s *Styles) CommandName(v string) string {
	return string(s.CommandNameCode) + v + string(s.DefaultCode)
}
// ErrorAlt wraps v with the alternate error style (used for highlighted values in errors).
func (s *Styles) ErrorAlt(v string) string {
	return string(s.ErrorAltCode) + v + string(s.DefaultCode)
}
// Error wraps v with the error style.
func (s *Styles) Error(v string) string { return string(s.ErrorCode) + v + string(s.DefaultCode) }

// FlagAlias wraps v with the flag alias style.
func (s *Styles) FlagAlias(v string) string {
	return string(s.FlagAliasCode) + v + string(s.DefaultCode)
}
// FlagName wraps v with the flag name style.
func (s *Styles) FlagName(v string) string { return string(s.FlagNameCode) + v + string(s.DefaultCode) }

// Header wraps v with the header style.
func (s *Styles) Header(v string) string   { return string(s.HeaderCode) + v + string(s.DefaultCode) }

// Label wraps v with the label style.
func (s *Styles) Label(v string) string {
	return string(s.LabelCode) + v + string(s.DefaultCode)
}
// SectionName wraps v with the section name style.
func (s *Styles) SectionName(v string) string {
	return string(s.SectionNameCode) + v + string(s.DefaultCode)
}

// Printer is a convenience wrapper around [fmt.Fprint], [fmt.Fprintln], and
// [fmt.Fprintf] with a fixed [io.Writer].
type Printer struct {
	w io.Writer
}

// NewPrinter creates a [Printer] that writes to w.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

// Println writes values followed by a newline to the underlying writer.
func (p *Printer) Println(v ...any) (int, error) {
	return fmt.Fprintln(p.w, v...)
}

// Print writes values to the underlying writer without a trailing newline.
func (p *Printer) Print(v ...any) (int, error) {
	return fmt.Fprint(p.w, v...)
}

// Printf writes a formatted string to the underlying writer.
func (p *Printer) Printf(format string, v ...any) (int, error) {
	return fmt.Fprintf(p.w, format, v...)
}
