// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package worker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/pkg/triggers"
)

var _ = Describe("Worker", func() {

	Context("UnimplementedWorker", func() {
		uiWrkr := &UnimplementedWorker{}

		When("calling HandlesEvent", func() {
			It("should return false", func() {
				Expect(uiWrkr.HandlesEvent(&triggers.Event{})).To(BeFalse())
			})
		})

		When("calling HandleEvent", func() {
			It("should return an error", func() {
				err := uiWrkr.HandleEvent(&triggers.Event{})
				Expect(err).Should(HaveOccurred())
			})
		})

		When("calling HandlesHttpRequest", func() {
			It("should return false", func() {
				Expect(uiWrkr.HandlesHttpRequest(&triggers.HttpRequest{})).To(BeFalse())
			})
		})

		When("calling HandleHttpRequest", func() {
			It("should return an error", func() {
				_, err := uiWrkr.HandleHttpRequest(&triggers.HttpRequest{})
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
