package membrane_test

import (
	"fmt"
	"plugin"
	"strings"

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

type MockPlugin struct {
	membrane.PluginIface
	SymbolMap map[string]interface{}
}

func (p *MockPlugin) Lookup(name string) (plugin.Symbol, error) {
	if symbol, ok := p.SymbolMap[name]; ok {
		return symbol, nil
	}

	return nil, fmt.Errorf("No such symbol found")
}

var _ = Describe("Membrane", func() {
	Context("Starting a new membrane, that tolerates missing services", func() {

		When("It is missing the gateway plugin", func() {
			mockPluginLoader := func(location string) (membrane.PluginIface, error) {
				return nil, fmt.Errorf("Failed to load plugin: %s", location)
			}

			membrane, _ := membrane.NewWithPluginLoader(&membrane.MembraneOptions{
				TolerateMissingServices: true,
			}, mockPluginLoader)
			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("The Gateway plugin is available but is missing a New() constructor", func() {
			mockPluginLoader := func(location string) (membrane.PluginIface, error) {
				if strings.Contains(location, "gateway") {
					return &MockPlugin{
						SymbolMap: map[string]interface{}{
							// Create a new Gateway
							"FakeMethod": func() (sdk.GatewayPlugin, error) {
								return nil, fmt.Errorf("There was an error creating the gateway")
							},
						},
					}, nil
				}

				return nil, fmt.Errorf("Failed to load plugin: %s", location)
			}

			membrane, _ := membrane.NewWithPluginLoader(&membrane.MembraneOptions{
				GatewayPluginFile:       "gateway.so",
				TolerateMissingServices: true,
			}, mockPluginLoader)
			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("The Gateway plugin is available but implements the wrong interface", func() {
			mockPluginLoader := func(location string) (membrane.PluginIface, error) {
				if strings.Contains(location, "gateway") {
					return &MockPlugin{
						SymbolMap: map[string]interface{}{
							// Create a new Gateway
							"New": func() (string, error) {
								return "Testing", nil
							},
						},
					}, nil
				}

				return nil, fmt.Errorf("Failed to load plugin: %s", location)
			}

			membrane, _ := membrane.NewWithPluginLoader(&membrane.MembraneOptions{
				GatewayPluginFile:       "gateway.so",
				TolerateMissingServices: true,
			}, mockPluginLoader)
			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("The Gateway plugin is available but returns an error", func() {
			mockPluginLoader := func(location string) (membrane.PluginIface, error) {
				if strings.Contains(location, "gateway") {
					return &MockPlugin{
						SymbolMap: map[string]interface{}{
							// Create a new Gateway
							"New": func() (sdk.GatewayPlugin, error) {
								return nil, fmt.Errorf("There was an error creating the gateway")
							},
						},
					}, nil
				}

				return nil, fmt.Errorf("Failed to load plugin: %s", location)
			}

			membrane, _ := membrane.NewWithPluginLoader(&membrane.MembraneOptions{
				GatewayPluginFile:       "gateway.so",
				TolerateMissingServices: true,
			}, mockPluginLoader)
			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("The Gateway plugin is available and working", func() {
			mockGateway := &MockGateway{}
			mockPluginLoader := func(location string) (membrane.PluginIface, error) {
				if strings.Contains(location, "gateway") {
					return &MockPlugin{
						SymbolMap: map[string]interface{}{
							// Create a new Gateway
							"New": func() (sdk.GatewayPlugin, error) {
								return mockGateway, nil
							},
						},
					}, nil
				}

				return nil, fmt.Errorf("Failed to load plugin: %s", location)
			}

			membrane, _ := membrane.NewWithPluginLoader(&membrane.MembraneOptions{
				GatewayPluginFile:       "gateway.so",
				TolerateMissingServices: true,
			}, mockPluginLoader)
			It("Start should not Panic", func() {
				Expect(membrane.Start).ToNot(Panic())
			})

			It("Mock Gateways start method should have been called", func() {
				Expect(mockGateway.started).To(BeTrue())
			})
		})
	})

	Context("Starting a new membrane, that does not tolerate missing services", func() {
		When("It is missing the eventing plugin", func() {
			mockPluginLoader := func(location string) (membrane.PluginIface, error) {
				return nil, fmt.Errorf("Failed to load plugin: %s", location)
			}

			membrane, _ := membrane.NewWithPluginLoader(&membrane.MembraneOptions{
				TolerateMissingServices: false,
			}, mockPluginLoader)
			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("It is missing the documents plugin", func() {
			mockEventingServer := &MockEventingServer{}
			mockPluginLoader := func(location string) (membrane.PluginIface, error) {
				if strings.Contains(location, "eventing") {
					return &MockPlugin{
						SymbolMap: map[string]interface{}{
							// Create a new Gateway
							"New": func() (sdk.EventingPlugin, error) {
								return mockEventingServer, nil
							},
						},
					}, nil
				}

				return nil, fmt.Errorf("Failed to load plugin: %s", location)
			}

			membrane, _ := membrane.NewWithPluginLoader(&membrane.MembraneOptions{
				EventingPluginFile:      "eventing.so",
				TolerateMissingServices: false,
			}, mockPluginLoader)

			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})

		When("It is missing the storage plugin", func() {
			mockEventingServer := &MockEventingServer{}
			mockDocumentsServer := &MockDocumentsServer{}
			mockPluginLoader := func(location string) (membrane.PluginIface, error) {
				if strings.Contains(location, "eventing") {
					return &MockPlugin{
						SymbolMap: map[string]interface{}{
							// Create a new Gateway
							"New": func() (sdk.EventingPlugin, error) {
								return mockEventingServer, nil
							},
						},
					}, nil
				}

				if strings.Contains(location, "documents") {
					return &MockPlugin{
						SymbolMap: map[string]interface{}{
							// Create a new Gateway
							"New": func() (sdk.DocumentsPlugin, error) {
								return mockDocumentsServer, nil
							},
						},
					}, nil
				}

				return nil, fmt.Errorf("Failed to load plugin: %s", location)
			}

			membrane, _ := membrane.NewWithPluginLoader(&membrane.MembraneOptions{
				EventingPluginFile:      "eventing.so",
				DocumentsPluginFile:     "documents.so",
				TolerateMissingServices: false,
			}, mockPluginLoader)

			It("Start should Panic", func() {
				Expect(membrane.Start).To(Panic())
			})
		})
	})
})
