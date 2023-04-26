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

package errors_test

import (
	"fmt"

	"github.com/kubespress/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func functionReturningError() error {
	return errors.New("")
}

func ExampleSet() {
	// UserFacingMessage is a type containing a user facing message
	type UserFacingMessage string

	// Enrich error using Set method
	err := errors.Enrich(functionReturningError(), errors.Set[UserFacingMessage]("Well this sucks."))

	// Print the error
	fmt.Println(errors.Get[UserFacingMessage](err, "Internal server error"))
	// Output: Well this sucks.
}

func ExampleGet() {
	// UserFacingMessage is a type containing a user facing message
	type UserFacingMessage string

	// Call function returning error
	err := functionReturningError()
	if err != nil {
		// Get the user facing message
		msg := errors.Get[UserFacingMessage](err, "Internal server error")

		// Print the error message
		fmt.Println(msg)
		// Output: Internal server error
	}
}

func ExampleCheck() {
	// Temporary indicates the error is Temporary
	type Temporary bool

	// Call function returning error
	err := functionReturningError()
	if err != nil {
		// If the error is Temporary then do nothing
		if errors.Check[Temporary](err) {
			return
		}

		// Print the error message
		fmt.Println(err)
	}
}

var _ = Describe("New", func() {
	var err error

	BeforeEach(func() {
		err = errors.New("test message 01")
	})

	It("should return a new error", func() {
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("test message 01"))
	})
})

var _ = Describe("Errorf", func() {
	var err error

	BeforeEach(func() {
		err = errors.Errorf("test message %02d", 2)
	})

	It("should return a new error", func() {
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("test message 02"))
	})
})

var _ = Describe("Enrich", func() {
	var err error

	Context("nil error", func() {
		BeforeEach(func() {
			err = errors.Enrich(nil)
		})

		It("should return nil", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("non nil error", func() {
		BeforeEach(func() {
			err = errors.New("test message 03")
		})

		Context("when enricher drops the error", func() {
			BeforeEach(func() {
				err = errors.Enrich(err, func(err error) error {
					return nil
				})
			})

			It("should return nil", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when enricher replaces the error", func() {
			BeforeEach(func() {
				err = errors.Enrich(err, func(err error) error {
					return errors.New("test message 04")
				})
			})

			It("should return replaced error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("test message 04"))
			})
		})

		Context("when context is added to the error", func() {
			type UnsetContextString string
			type UnsetContextInt int
			type UnsetContextBool bool
			type ContextString string
			type ContextInt int
			type ContextBool bool

			BeforeEach(func() {
				err = errors.Enrich(err,
					errors.Set[ContextString]("additional context"),
					errors.Set[ContextInt](123),
					errors.Set[ContextBool](true),
				)
			})

			It("should be retrievable from the error", func() {
				Expect(errors.Get[ContextString](err, "default value")).To(Equal(ContextString("additional context")))
				Expect(errors.Get[ContextInt](err, 456)).To(Equal(ContextInt(123)))
				Expect(errors.Check[ContextBool](err)).To(BeTrue())
				Expect(errors.Get[UnsetContextString](err, "default value")).To(Equal(UnsetContextString("default value")))
				Expect(errors.Get[UnsetContextInt](err, 456)).To(Equal(UnsetContextInt(456)))
				Expect(errors.Check[UnsetContextBool](err)).To(BeFalse())
				Expect(fmt.Sprintf("%+v", err)).To(Equal("test message 03"))
				Expect(fmt.Sprintf("%v", err)).To(Equal("test message 03"))
				Expect(fmt.Sprintf("%s", err)).To(Equal("test message 03"))
				Expect(fmt.Sprintf("%q", err)).To(Equal(`"test message 03"`))
			})
		})

		Context("when error is wrapped", func() {
			BeforeEach(func() {
				err = errors.Enrich(err,
					errors.Wrapf("message prefix"),
				)
			})

			It("should prefix the error message", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("message prefix: test message 03"))
				Expect(fmt.Sprintf("%+v", err)).To(Equal("message prefix: test message 03"))
				Expect(fmt.Sprintf("%v", err)).To(Equal("message prefix: test message 03"))
				Expect(fmt.Sprintf("%s", err)).To(Equal("message prefix: test message 03"))
				Expect(fmt.Sprintf("%q", err)).To(Equal(`"message prefix: test message 03"`))
			})
		})
	})
})
