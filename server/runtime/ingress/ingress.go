package ingress

import (
	"fmt"

	"github.com/nitrictech/nitric/server/runtime"
)

type IngressPlugin interface {
	Start(Proxy) error
}

var ingressPlugin IngressPlugin = nil

// Register a new instance of a storage plugin
func Register(constructor runtime.PluginConstructor[IngressPlugin]) error {
	if ingressPlugin != nil {
		return fmt.Errorf("ingress plugin already registered")
	}

	ingressPlugin = constructor()

	return nil
}
