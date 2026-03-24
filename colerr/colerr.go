package colerr

import (
	"errors"
	"fmt"
	"io"

	"go.bbkane.com/warg/styles"
)

// Colerr provides a way to create errors that can be printed with color styles. The Stacktrace function can be used to print the error and it wrapped errors with the appropriate styles.
//
// Types the implement ColerError should also implement Error (and optionally Unwrap) so that they can be used as normal errors in the repl
type ColorError interface {
	// ColorError returns a string formatted with the style.
	//
	// If the error is a wrapper around another error, the ColorError method should only format the message of this error, not the wrapped error. The Stacktrace function will handle printing the wrapped error separately.
	ColorError(s *styles.Styles) string
}

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

func NewWrapped(err error, msg string) Wrapped {
	return Wrapped{
		err: err,
		msg: msg,
	}
}

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

func NewWrappedf(err error, format string, args ...string) Wrappedf {
	return Wrappedf{
		err:  err,
		msg:  format,
		args: args,
	}
}
