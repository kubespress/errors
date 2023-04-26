/*
Copyright 2023 Kubespress Authors.
*/

package errors

import "strings"

type errorAggregate struct {
	errs []error
}

// ErrorList is a list of errors, it can be used as a utility to aggregate
// errors by appending to the slice
type ErrorList []error

// Error returns the aggregated error, if the list is empty, nil is returned
func (err ErrorList) Error() error {
	return Aggregate(err...)
}

// Aggregate returns an error wrapping multiple other errors. If no other errors
// are passed in this method returns nil.
func Aggregate(errs ...error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errorAggregate{errs: errs}
	}
}

func (e errorAggregate) Error() string {
	// Track seen errors and their message
	seenerrs := map[string]struct{}{}
	messages := make([]string, 0, len(e.errs))

	// Action depends on the length
	switch len(e.errs) {
	case 0:
		return ""
	case 1:
		return e.errs[0].Error()
	default:
		e.visit(func(err error) {
			// Check to see if we have already seen this error
			msg := err.Error()
			if _, seen := seenerrs[msg]; seen {
				return
			}

			// Insert this error into the map and add its contents to the message
			seenerrs[msg] = struct{}{}
			messages = append(messages, msg)
		})
	}

	// Only one real error, just return the message
	if len(messages) == 1 {
		return messages[0]
	}

	// Return the joined messages
	return "[" + strings.Join(messages, ", ") + "]"
}

func (e errorAggregate) visit(fn func(error)) {
	type aggregate interface {
		Errors() []error
	}

	// Loop over errors
	for _, err := range e.errs {
		switch err := err.(type) {
		// If the error is another errorAggregate, recurse into its visit method
		case errorAggregate:
			err.visit(fn)

		// If the error has an Errors() method, convert to a errorAggregate and
		// visit all its errors. This allows it to work with similar libraries.
		case aggregate:
			errorAggregate{errs: err.Errors()}.visit(fn)

		// Call the visit function on the error
		default:
			fn(err)
		}
	}
}

func (e errorAggregate) Unwrap() []error {
	return e.errs
}

func (e errorAggregate) Errors() []error {
	return e.Unwrap()
}
