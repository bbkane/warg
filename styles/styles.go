package styles

import (
	"fmt"
	"io"

	"go.bbkane.com/gocolor"
)

// Replace the fmt* functions I'm using for coloring
// Create a styles struct with public fileds for coloring sections, headers, etc

type Styles struct {
	CommandNameCode gocolor.Code
	FlagAliasCode   gocolor.Code
	FlagNameCode    gocolor.Code
	HeaderCode      gocolor.Code
	LabelCode       gocolor.Code
	SectionNameCode gocolor.Code

	// DefaultCode should probably always be to gocolor.Default or gocolor.Empty
	DefaultCode gocolor.Code
}

func NewEmptyStyles() Styles {
	return Styles{
		CommandNameCode: gocolor.Empty,
		FlagAliasCode:   gocolor.Empty,
		FlagNameCode:    gocolor.Empty,
		HeaderCode:      gocolor.Empty,
		LabelCode:       gocolor.Empty,
		SectionNameCode: gocolor.Empty,
		DefaultCode:     gocolor.Empty,
	}
}

func NewEnabledStyles() Styles {
	return Styles{
		CommandNameCode: gocolor.Bold + gocolor.FgGreen,
		FlagAliasCode:   gocolor.Bold + gocolor.FgYellow,
		FlagNameCode:    gocolor.Bold + gocolor.FgYellow,
		HeaderCode:      gocolor.Bold + gocolor.Underline,
		LabelCode:       gocolor.Bold,
		SectionNameCode: gocolor.Bold + gocolor.FgCyan,
		DefaultCode:     gocolor.Default,
	}
}

func (s *Styles) CommandName(v string) string {
	return string(s.CommandNameCode) + v + string(s.DefaultCode)
}
func (s *Styles) FlagAlias(v string) string {
	return string(s.FlagAliasCode) + v + string(s.DefaultCode)
}
func (s *Styles) FlagName(v string) string { return string(s.FlagNameCode) + v + string(s.DefaultCode) }
func (s *Styles) Header(v string) string   { return string(s.HeaderCode) + v + string(s.DefaultCode) }
func (s *Styles) Label(v string) string {
	return string(s.LabelCode) + v + string(s.DefaultCode)
}
func (s *Styles) SectionName(v string) string {
	return string(s.SectionNameCode) + v + string(s.DefaultCode)
}

// Printer is a small convenience wrapper around fmt.FprintXX and an io.Writer
type Printer struct {
	w io.Writer
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

func (p *Printer) Println(v ...any) (int, error) {
	return fmt.Fprintln(p.w, v...)
}

func (p *Printer) Print(v ...any) (int, error) {
	return fmt.Fprint(p.w, v...)
}

func (p *Printer) Printf(format string, v ...any) (int, error) {
	return fmt.Fprintf(p.w, format, v...)
}
