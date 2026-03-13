package colerr

import (
	"errors"
	"fmt"
	"io"

	"go.bbkane.com/warg/styles"
)

func Stacktrace(w io.Writer, style *styles.Styles, err error) {
	p := styles.NewPrinter(w)

	// If the error implements ColorError, use that instead of the default Error string
	type colorErr interface {
		ColorError(s *styles.Styles) string
	}

	if ce, ok := err.(colorErr); ok {
		p.Print(ce.ColorError(style))
	} else {
		p.Print(style.Error(err.Error()))
	}

	under := errors.Unwrap(err)
	if under != nil {
		p.Println()
		p.Println()
		Stacktrace(w, style, under)
	}
}

type Wrapped struct {
	err error
	msg string
}

func (w Wrapped) Error() string {
	return w.msg
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
	return fmt.Sprintf(w.msg, args...)
}

func (w Wrappedf) ColorError(s *styles.Styles) string {
	var args []any
	for _, a := range w.args {
		args = append(args, s.ErrorAlt(a))
	}
	return s.Error(fmt.Sprintf(w.msg, args...))
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
