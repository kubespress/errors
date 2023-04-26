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
	"github.com/kubespress/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aggregate", func() {
	var err error
	var errs errors.ErrorList

	Context("with no errors", func() {
		BeforeEach(func() {
			err = errs.Error()
		})

		It("should return nil", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("with one error", func() {
		var err1 = errors.New("example error 07")

		BeforeEach(func() {
			errs = errors.ErrorList{err1}
			err = errs.Error()
		})

		It("should return the single error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(err1))
		})
	})

	Context("with multiple errors", func() {
		var err1 = errors.New("example error 07")
		var err2 = errors.New("example error 08")

		BeforeEach(func() {
			errs = errors.ErrorList{err1, err2}
			err = errs.Error()
		})

		It("should aggregated error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("[example error 07, example error 08]"))
			Expect(errors.Is(err, err1)).To(BeTrue())
			Expect(errors.Is(err, err2)).To(BeTrue())
		})
	})
})
