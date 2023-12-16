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
				Expect(fmt.Sprintf("%+v", err)).To(MatchRegexp("^test message 05\ngithub\\.com/kubespress/errors_test\\.glob\\.\\.func.*\n\t.*/stack_test.go:.*"))
				Expect(fmt.Sprintf("%v", err)).To(Equal("test message 05"))
				Expect(fmt.Sprintf("%s", err)).To(Equal("test message 05"))
				Expect(fmt.Sprintf("%q", err)).To(Equal(`"test message 05"`))
			})
		})
	})
})
