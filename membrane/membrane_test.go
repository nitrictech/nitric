package membrane_test

import (
	"github.com/nitric-dev/membrane/membrane"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockEventingServer struct {
	sdk.UnimplementedEventingPlugin
}

type MockStorageServer struct {
	sdk.UnimplementedStoragePlugin
}

type MockDocumentsServer struct {
	sdk.UnimplementedDocumentsPlugin
}

type MockGateway struct {
	sdk.UnimplementedGatewayPlugin
	started bool
}

func (gw *MockGateway) Start(handler sdk.GatewayHandler) error {
	// Spy on the mock gateway
	gw.started = true
	return nil
}

var _ = Describe("Membrane", func() {
	Context("Starting a new membrane, that tolerates missing services", func() {

		When("It is missing the gateway plugin", func() {
			membrane, _ := membrane.New(&membrane.MembraneOptions{
				TolerateMissingServices: true,
				SuppressLogs:            true,
			})
			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("The Gateway plugin is available and working", func() {
			mockGateway := &MockGateway{}

			membrane, _ := membrane.New(&membrane.MembraneOptions{
				GatewayPlugin:           mockGateway,
				SuppressLogs:            true,
				TolerateMissingServices: true,
			})

			It("Start should not Panic", func() {
				Expect(membrane.Start).ToNot(Panic())
			})

			It("Mock Gateways start method should have been called", func() {
				Expect(mockGateway.started).To(BeTrue())
			})
		})
	})

	Context("Starting a new membrane, that does not tolerate missing services", func() {
		mockGateway := &MockGateway{}
		When("It is missing the eventing plugin", func() {
			membrane, _ := membrane.New(&membrane.MembraneOptions{
				TolerateMissingServices: false,
				SuppressLogs:            true,
				GatewayPlugin:           mockGateway,
			})
			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("It is missing the documents plugin", func() {
			mockEventingServer := &MockEventingServer{}

			membrane, _ := membrane.New(&membrane.MembraneOptions{
				EventingPlugin:          mockEventingServer,
				GatewayPlugin:           mockGateway,
				SuppressLogs:            true,
				TolerateMissingServices: false,
			})

			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("It is missing the storage plugin", func() {
			mockEventingServer := &MockEventingServer{}
			mockDocumentsServer := &MockDocumentsServer{}

			membrane, _ := membrane.New(&membrane.MembraneOptions{
				EventingPlugin:          mockEventingServer,
				DocumentsPlugin:         mockDocumentsServer,
				GatewayPlugin:           mockGateway,
				SuppressLogs:            true,
				TolerateMissingServices: false,
			})

			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})
	})
})
