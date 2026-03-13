package colerr

import (
	"errors"
	"fmt"
	"io"

	"go.bbkane.com/warg/styles"
)

func Stacktrace(w io.Writer, style styles.Styles, err error) {
	p := styles.NewPrinter(w)
	p.Println(style.Error(err.Error()))

	under := errors.Unwrap(err)
	if under != nil {
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

func NewWrappedf(err error, format string, args ...any) Wrapped {
	return Wrapped{
		err: err,
		msg: fmt.Sprintf(format, args...),
	}
}
