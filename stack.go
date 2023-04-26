/*
Copyright 2023 Kubespress Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
