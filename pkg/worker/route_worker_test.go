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
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock "github.com/nitrictech/nitric/mocks/worker"
	"github.com/nitrictech/nitric/pkg/triggers"
)

var _ = Describe("RouteWorker", func() {

	Context("Http", func() {
		rWrkr := &RouteWorker{
			methods: []string{"GET"},
			path:    "/test/:param",
		}

		When("calling HandlesHttpRequest with bad path", func() {
			It("should return false", func() {
				Expect(rWrkr.HandlesHttpRequest(&triggers.HttpRequest{
					Method: "GET",
					Path:   "/test/",
				})).To(BeFalse())
			})
		})

		When("calling HandlesHttpRequest with bad method", func() {
			It("should return false", func() {
				Expect(rWrkr.HandlesHttpRequest(&triggers.HttpRequest{
					Method: "POST",
					Path:   "/test/test",
				})).To(BeFalse())
			})
		})

		When("calling HandlesHttpRequest with matching path and method", func() {
			It("should return true", func() {
				Expect(rWrkr.HandlesHttpRequest(&triggers.HttpRequest{
					Method: "GET",
					Path:   "/test/test",
				})).To(BeTrue())
			})
		})

		When("calling HandleHttpRequest", func() {
			It("should call the base grpc workers HandleEvent with augmented trigger", func() {
				ctrl := gomock.NewController(GinkgoT())
				mGrpc := mock.NewMockGrpcWorker(ctrl)

				By("calling the base grpc handler HandleEvent method")
				mGrpc.EXPECT().HandleHttpRequest(&triggers.HttpRequest{
					Method: "GET",
					Path:   "/test/name",
					Params: map[string]string{
						"param": "name",
					},
				}).Times(1)

				subWrkr := &RouteWorker{
					methods:    []string{"GET"},
					path:       "/test/:param",
					GrpcWorker: mGrpc,
				}

				_, err := subWrkr.HandleHttpRequest(&triggers.HttpRequest{
					Method: "GET",
					Path:   "/test/name",
				})

				Expect(err).ShouldNot(HaveOccurred())
				ctrl.Finish()
			})
		})
	})

	Context("Event", func() {

		When("calling HandlesEvent", func() {
			rWrkr := &RouteWorker{}

			It("should return false", func() {
				Expect(rWrkr.HandlesEvent(&triggers.Event{})).To(BeFalse())
			})
		})

		When("calling HandleEvent", func() {
			subWrkr := &RouteWorker{}

			It("should return an error", func() {

				err := subWrkr.HandleEvent(&triggers.Event{})

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
