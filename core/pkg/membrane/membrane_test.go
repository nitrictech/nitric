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

package membrane_test

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	mock_gateway "github.com/nitrictech/nitric/core/mocks/gateway"
	"github.com/nitrictech/nitric/core/pkg/membrane"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Membrane", func() {
	Context("Starting the server", func() {
		When("The Gateway plugin is available and working", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockGateway := mock_gateway.NewMockGatewayService(ctrl)

			os.Args = []string{}
			membrane, _ := membrane.New(
				membrane.WithMinWorkers(0),
				membrane.WithGatewayPlugin(mockGateway),
			)

			It("Should successfully start the membrane", func() {
				By("starting the gateway plugin")
				mockGateway.EXPECT().Start(gomock.Any()).Times(1).Return(nil)

				_ = membrane.Start()
			})
		})

		When("The configured service port is already consumed", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockGateway := mock_gateway.NewMockGatewayService(ctrl)
			mockGateway.EXPECT().Start(gomock.Any()).AnyTimes().Return(nil)
			var lis net.Listener

			membrane, _ := membrane.New(membrane.WithMinWorkers(0), membrane.WithGatewayPlugin(mockGateway), membrane.WithServiceAddress("localhost:9005"))

			BeforeEach(func() {
				lis, _ = net.Listen("tcp", "localhost:9005")
			})

			AfterEach(func() {
				lis.Close()
			})

			It("Should return an error", func() {
				err := membrane.Start()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("could not listen"))
			})
		})
	})

	Context("Starting the child process", func() {
		BeforeEach(func() {
			os.Args = []string{}
		})

		var mb *membrane.Membrane
		When("The configured command exists", func() {
			BeforeEach(func() {
				ctrl := gomock.NewController(GinkgoT())
				mockGateway := mock_gateway.NewMockGatewayService(ctrl)
				mockGateway.EXPECT().Start(gomock.Any()).AnyTimes().Return(nil)
				mockGateway.EXPECT().Stop().AnyTimes().Return(nil)
				mb, _ = membrane.New(
					membrane.WithChildCommand([]string{"sleep", "5"}),
					membrane.WithGatewayPlugin(mockGateway),
					membrane.WithServiceAddress(fmt.Sprintf(":%d", 9001)),
					membrane.WithChildTimeoutSeconds(1),
				)
			})

			AfterEach(func() {
				mb.Stop()
			})

			When("There is a worker available in the pool", func() {
				BeforeEach(func() {
					go (func() {
						_ = http.ListenAndServe("localhost:8081", nil)
					})()
				})

				It("Should wait for the service to start", func() {
					_ = mb.Start()
				})
			})
		})

		When("The configured command does not exist", func() {
			BeforeEach(func() {
				ctrl := gomock.NewController(GinkgoT())
				mockGateway := mock_gateway.NewMockGatewayService(ctrl)

				mb, _ = membrane.New(
					membrane.WithChildCommand([]string{"fakecommand"}),
					membrane.WithGatewayPlugin(mockGateway),
				)
			})

			It("Should return an error", func() {
				err := mb.Start()
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
