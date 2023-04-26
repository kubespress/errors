/*
Copyright 2023 Kubespress Authors.
*/

package errors

import (
	"fmt"
	"io"
	"runtime"
)

type errWithStack struct {
	stack []uintptr
	err   error
}

func (e errWithStack) Error() string { return e.err.Error() }
func (e errWithStack) Unwrap() error { return e.err }

func (w errWithStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Unwrap())
			for _, pc := range w.stack {
				fmt.Fprint(s, "\n", w.line(pc-1))
			}

			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

func (w *errWithStack) line(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown\n\tunknown:0"
	}

	file, line := fn.FileLine(pc)
	return fmt.Sprintf("%s\n\t%s:%d", fn.Name(), file, line)
}

// WithStack will enrich the message with a stacktrace
func WithStack() Enricher {
	return WithStackN(1)
}

// WithStackN will enrich the message with a stacktrace, skipping a specified
// number of calls.
func WithStackN(skip int) Enricher {
	return func(err error) error {
		const depth = 32
		var pcs [depth]uintptr
		n := runtime.Callers(skip+2, pcs[:])

		return errWithStack{
			stack: pcs[0:n],
			err:   err,
		}
	}
}
