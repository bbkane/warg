package main

import (
	"errors"
	"fmt"
	"io"
	"os"

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

func main() {
	s := styles.NewEnabledStyles()

	fmt.Println("# just one one err")
	Stacktrace(os.Stderr, s, errors.New("this is an error"))

	fmt.Println("# wrapped fmt.Errorf")
	Stacktrace(os.Stderr, s, fmt.Errorf("this is a wrapped error: %w", errors.New("this is the inner error")))

	fmt.Println("# wrapped custom error")
	Stacktrace(os.Stderr, s, NewWrapped(errors.New("wrapped err"), "wrapper msg"))
}
