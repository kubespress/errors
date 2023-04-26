/*
Copyright 2023 Kubespress Authors.
*/

package errors_test

import (
	"fmt"

	"github.com/kubespress/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Enrich", func() {
	var err error

	Context("non nil error", func() {
		BeforeEach(func() {
			err = errors.New("test message 05")
		})

		Context("with stack", func() {
			BeforeEach(func() {
				err = errors.Enrich(err, errors.WithStack())
			})

			It("should return error with stack", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("test message 05"))
				Expect(fmt.Sprintf("%+v", err)).To(MatchRegexp("^test message 05\ngithub\\.com/kubespress/errors_test\\.glob\\.\\.func.*\n\t/workspaces/errors/stack_test.go:.*"))
				Expect(fmt.Sprintf("%v", err)).To(Equal("test message 05"))
				Expect(fmt.Sprintf("%s", err)).To(Equal("test message 05"))
				Expect(fmt.Sprintf("%q", err)).To(Equal(`"test message 05"`))
			})
		})
	})
})
