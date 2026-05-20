package colerr

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"go.bbkane.com/warg/styles"
)

// Package colerr provides colored error types that support styled terminal output.
// Types implementing [ColorError] can be printed with [Stacktrace] to display a
// chain of wrapped errors with ANSI color styling.

// ColorError is implemented by errors that support colored terminal output.
// Implementations should also implement the standard error interface
// (and optionally Unwrap) so they work as normal Go errors.
type ColorError interface {
	// ColorError returns a styled string for this error's message only (not wrapped errors, unlike the standard Go convention).
	// [Stacktrace] handles printing wrapped errors separately.
	ColorError(s *styles.Styles) string
}

// Stacktrace prints err and all its wrapped errors to w using the given styles.
// Each [ColorError] in the chain is printed on its own line; the final non-ColorError
// is printed with the error style. Multi-errors (Unwrap() []error) are indented.
func Stacktrace(w io.Writer, style *styles.Styles, err error) {
	p := styles.NewPrinter(w)

	// loop through all color errors, printing them until we find one that's not
	for {
		ce, isColerError := err.(ColorError)
		if !isColerError {
			break
		}
		p.Println(ce.ColorError(style))

		err = errors.Unwrap(err)
		if err == nil {
			return
		}

		// if the wrapped error is a slice, consider this the last error and print each error in the slice in an indented new line, then stop
		if unwrapped, ok := err.(interface{ Unwrap() []error }); ok {
			p.Println()
			for _, e := range unwrapped.Unwrap() {
				// Assume these are simple
				p.Println("  ", style.Error(e.Error()))
			}
			return
		}

		p.Println()

	}

	// print the last error, which is not a color error
	p.Println(style.Error(err.Error()))
}

// Wrapped is an error that wraps another error with an additional message.
// Implements [ColorError], error, and Unwrap.
type Wrapped struct {
	err error
	msg string
}

func (w Wrapped) ColorError(s *styles.Styles) string {
	return s.Error(w.msg)
}

func (w Wrapped) Error() string {
	return w.msg + ": " + w.err.Error()
}

func (w Wrapped) Unwrap() error {
	return w.err
}

// NewWrapped creates a [Wrapped] error wrapping err with the given message.
// err may be nil for root-cause errors without an underlying cause.
func NewWrapped(err error, msg string) Wrapped {
	return Wrapped{
		err: err,
		msg: msg,
	}
}

// Wrappedf is an error that wraps another error with a formatted message.
// Format arguments are styled with ErrorAltCode when printed via [Stacktrace].
// Implements [ColorError], error, and Unwrap.
type Wrappedf struct {
	err  error
	msg  string
	args []string
}

func (w Wrappedf) Error() string {
	var args []any
	for _, a := range w.args {
		args = append(args, a)
	}
	return fmt.Sprintf(w.msg, args...) + ": " + w.err.Error()
}

func (w Wrappedf) ColorError(s *styles.Styles) string {
	var args []any
	for _, a := range w.args {
		// manually apply ErrorAltCode since we need to reset it to the ErrorCode instead of the default mid-string
		a := string(s.ErrorAltCode) + a + string(s.ErrorCode)
		args = append(args, a)
	}
	return fmt.Sprintf(s.Error(w.msg), args...)
}

func (w Wrappedf) Unwrap() error {
	return w.err
}

// NewWrappedf creates a [Wrappedf] error wrapping err with a fmt.Sprintf-style message.
// err may be nil for root-cause errors without an underlying cause.
func NewWrappedf(err error, format string, args ...string) Wrappedf {
	return Wrappedf{
		err:  err,
		msg:  format,
		args: args,
	}
}

// ArgChoiceError is returned when a parsed argument does not match any valid choice.
// Implements [ColorError] for styled output listing the valid options.
type ArgChoiceError struct {
	// Message describes what was expected (e.g., "expecting section or command")
	Message string
	// Arg is the actual argument that was received
	Arg string
	// Choices lists the valid options
	Choices []string
}

func (e ArgChoiceError) Error() string {
	return fmt.Sprintf("%s, got %s. Choices: %v", e.Message, e.Arg, e.Choices)
}

func (e ArgChoiceError) ColorError(s *styles.Styles) string {
	var buf strings.Builder
	buf.WriteString(e.Message + "\n")
	buf.WriteString("Got: " + string(s.ErrorAltCode) + e.Arg + string(s.ErrorCode) + "\n")
	buf.WriteString("Choices:\n")
	buf.WriteString(string(s.ErrorAltCode))
	for _, c := range e.Choices {
		buf.WriteString("  " + c + "\n")
	}
	buf.WriteString(string(s.ErrorCode))

	return s.Error(buf.String())
}
