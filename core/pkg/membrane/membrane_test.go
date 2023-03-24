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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_gateway "github.com/nitrictech/nitric/core/mocks/gateway"
	mock_pool "github.com/nitrictech/nitric/core/mocks/pool"
	"github.com/nitrictech/nitric/core/pkg/membrane"
	"github.com/nitrictech/nitric/core/pkg/plugins/document"
	"github.com/nitrictech/nitric/core/pkg/plugins/events"
	"github.com/nitrictech/nitric/core/pkg/plugins/queue"
	"github.com/nitrictech/nitric/core/pkg/plugins/storage"
)

type MockDocumentServer struct {
	document.UnimplementedDocumentPlugin
}

type MockeventsServer struct {
	events.UnimplementedeventsPlugin
}

type MockStorageServiceServer struct {
	storage.UnimplementedStoragePlugin
}

type MockQueueServiceServer struct {
	queue.UnimplementedQueuePlugin
}

var _ = Describe("Membrane", func() {
	ctrl := gomock.NewController(GinkgoT())
	mockPool := mock_pool.NewMockWorkerPool(ctrl)

	BeforeSuite(func() {
		mockPool.EXPECT().WaitForMinimumWorkers(gomock.Any()).AnyTimes().Return(nil)
		mockPool.EXPECT().Monitor().AnyTimes().Return(nil)
		mockPool.EXPECT().GetWorkerCount().AnyTimes().Return(1)
		mockPool.EXPECT().GetMaxWorkers().AnyTimes().Return(100)
		os.Args = []string{}
	})

	Context("New", func() {
		Context("Tolerate Missing Services is enabled", func() {
			When("The gateway plugin is missing", func() {
				It("Should still fail to create", func() {
					m, err := membrane.New(&membrane.MembraneOptions{
						SuppressLogs:            true,
						TolerateMissingServices: true,
					})
					Expect(err).Should(HaveOccurred())
					Expect(m).To(BeNil())
				})
			})

			When("The gateway plugin is present", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockGateway := mock_gateway.NewMockGatewayService(ctrl)

				mbraneOpts := membrane.MembraneOptions{
					SuppressLogs:            true,
					GatewayPlugin:           mockGateway,
					TolerateMissingServices: true,
					Pool:                    mockPool,
				}
				It("Should successfully create the membrane server", func() {
					m, err := membrane.New(&mbraneOpts)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(m).ToNot(BeNil())
				})
			})
		})

		Context("Tolerate Missing Services is disabled", func() {
			When("Only the gateway plugin is present", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockGateway := mock_gateway.NewMockGatewayService(ctrl)

				mbraneOpts := membrane.MembraneOptions{
					TolerateMissingServices: false,
					SuppressLogs:            true,
					GatewayPlugin:           mockGateway,
					Pool:                    mockPool,
				}
				It("Should fail to create", func() {
					m, err := membrane.New(&mbraneOpts)
					Expect(err).Should(HaveOccurred())
					Expect(m).To(BeNil())
				})
			})

			When("All plugins are present", func() {
				mockDocumentServer := &MockDocumentServer{}
				mockeventsServer := &MockeventsServer{}
				mockStorageServiceServer := &MockStorageServiceServer{}
				mockQueueServiceServer := &MockQueueServiceServer{}
				ctrl := gomock.NewController(GinkgoT())
				mockGateway := mock_gateway.NewMockGatewayService(ctrl)

				mbraneOpts := membrane.MembraneOptions{
					TolerateMissingServices: false,
					SuppressLogs:            true,
					GatewayPlugin:           mockGateway,
					DocumentPlugin:          mockDocumentServer,
					EventsPlugin:            mockeventsServer,
					StoragePlugin:           mockStorageServiceServer,
					QueuePlugin:             mockQueueServiceServer,
					Pool:                    mockPool,
				}

				It("Should successfully create the membrane server", func() {
					m, err := membrane.New(&mbraneOpts)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(m).ToNot(BeNil())
				})
			})
		})
	})

	Context("Starting the server", func() {
		Context("That tolerates missing adapters", func() {
			When("The Gateway plugin is available and working", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockGateway := mock_gateway.NewMockGatewayService(ctrl)

				os.Args = []string{}
				membrane, _ := membrane.New(&membrane.MembraneOptions{
					GatewayPlugin:           mockGateway,
					SuppressLogs:            true,
					TolerateMissingServices: true,
					Pool:                    mockPool,
				})

				It("Should successfully start the membrane", func() {
					By("starting the gateway plugin")
					mockGateway.EXPECT().Start(mockPool).Times(1).Return(nil)

					// FIXME: Race condition causing inconsistent error here
					_ = membrane.Start()
				})
			})
		})

		When("The configured service port is already consumed", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockGateway := mock_gateway.NewMockGatewayService(ctrl)
			mockGateway.EXPECT().Start(gomock.Any()).AnyTimes().Return(nil)
			var lis net.Listener

			membrane, _ := membrane.New(&membrane.MembraneOptions{
				GatewayPlugin:           mockGateway,
				SuppressLogs:            true,
				TolerateMissingServices: true,
				ServiceAddress:          "localhost:9005",
				Pool:                    mockPool,
			})

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
				mb, _ = membrane.New(&membrane.MembraneOptions{
					ChildCommand:            []string{"sleep", "5"},
					GatewayPlugin:           mockGateway,
					ServiceAddress:          fmt.Sprintf(":%d", 9001),
					ChildTimeoutSeconds:     1,
					TolerateMissingServices: true,
					SuppressLogs:            true,
					Pool:                    mockPool,
				})
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
					// FIXME: Inconsistent error return with mocks
					_ = mb.Start()
				})
			})
		})

		When("The configured command does not exist", func() {
			BeforeEach(func() {
				ctrl := gomock.NewController(GinkgoT())
				mockGateway := mock_gateway.NewMockGatewayService(ctrl)

				mb, _ = membrane.New(&membrane.MembraneOptions{
					ChildAddress:            "localhost:808",
					ChildCommand:            []string{"fakecommand"},
					GatewayPlugin:           mockGateway,
					TolerateMissingServices: true,
					SuppressLogs:            true,
				})
			})

			It("Should return an error", func() {
				err := mb.Start()
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
