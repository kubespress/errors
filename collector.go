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

// ErrorCollector is a wrapper for ErrorList that provides methods for
// adding errors if they are not nil
type ErrorCollector struct {
	ErrorList
}

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{}
}

// CollectError will append the error to the ErrorList if it is not nil, it is
// designed to allow for a pattern where function call is passed right into
// the method, for example:
//
// collection.CollectError(someFunctionThatMayError())
func (e *ErrorCollector) CollectError(err error) {
	if err != nil {
		e.ErrorList = append(e.ErrorList, err)
	}
}

// TypedErrorCollector that allows you to pass through return values while
// keeping error
type TypedErrorCollector[T any] struct {
	collector *ErrorCollector
}

// ErrorCollectorForType returns a TypedErrorCollector that will collect errors
// onto a given ErrorCollector
func ErrorCollectorForType[T any](collector *ErrorCollector) TypedErrorCollector[T] {
	return TypedErrorCollector[T]{
		collector: collector,
	}
}

// CollectError will append the error to the ErrorList if it is not nil and
// return the specified value. It is designed to allow for a pattern where
// function call is passed right into the method, for example:
//
// value := collection.CollectError(someFunctionThatMayError())
func (e TypedErrorCollector[T]) CollectError(result T, err error) T {
	e.collector.CollectError(err)
	return result
}
