/*
Copyright 2023 Kubespress Authors.
*/

package errors

import (
	"errors"
	"fmt"
)

type errorString struct {
	msg string
}

func (e errorString) Error() string { return e.msg }

// New creates a new error
func New(msg string) error {
	return errorString{msg: msg}
}

// Errorf creates a new error using string formatting.
func Errorf(msg string, args ...any) error {
	return fmt.Errorf(msg, args...)
}

// Enricher is an error enrichment, it adds additional context to the provided
// error.
type Enricher func(error) error

// Enrich adds additional context to an error using the provided enrichment
// functions
func Enrich(err error, enrichments ...Enricher) error {
	// If error is nil, return nil
	if err == nil {
		return nil
	}

	// Loop over each enrichment function
	for _, e := range enrichments {
		// Enrich using the provided function
		err = e(err)

		// Enrichments can drop the error entirely, if that happens return nil
		if err == nil {
			return nil
		}
	}

	// Return enriched error
	return err
}

type enrichedError[T any] struct {
	enrichment T
	nested     error
}

func (err enrichedError[T]) Error() string { return err.nested.Error() }
func (err enrichedError[T]) Unwrap() error { return err.nested }
func (err enrichedError[T]) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", err.Unwrap())
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", err.Error())
	}
}

// Set enriches an error with a specific type.
func Set[T any](value T) Enricher {
	return func(err error) error {
		return enrichedError[T]{
			enrichment: value,
			nested:     err,
		}
	}
}

// Is returns true if the error contains the specified context and it is true.
func Is[T ~bool](err error) bool {
	// Check if error has enrichment
	var unwrapped enrichedError[T]
	if errors.As(err, &unwrapped) {
		// Cast enrichment to bool and return
		return bool(unwrapped.enrichment)
	}

	// Does not have enrichment
	return false
}

// Get returns the enriched context if it exists in the error, otherwise it
// returns the provided default value.
func Get[T any](err error, def T) T {
	// Check if error has enrichment
	var unwrapped enrichedError[T]
	if errors.As(err, &unwrapped) {
		return unwrapped.enrichment
	}

	// Does not have enrichment
	return def
}

type wrappedError struct {
	msg    string
	nested error
}

func (err wrappedError) Error() string { return err.msg + ": " + err.nested.Error() }
func (err wrappedError) Unwrap() error { return err.nested }
func (err wrappedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s: %+v", err.msg, err.Unwrap())
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", err.Error())
	}
}

// Wrap returns an enricher that prefixes an error message with a provided
// string for additional context
func Wrap(msg string) Enricher {
	return func(err error) error {
		return wrappedError{
			msg:    msg,
			nested: err,
		}
	}
}