package service

import (
	"fmt"

	"github.com/nitrictech/nitric/runtime/plugin"
)

type Service interface {
	Start(Proxy) error
}

var servicePlugin Service = nil

// Register a new instance of a storage plugin
func Register(constructor plugin.Constructor[Service]) error {
	var err error
	if servicePlugin != nil {
		return fmt.Errorf("ingress plugin already registered")
	}

	servicePlugin, err = constructor()
	if err != nil {
		return err
	}

	return nil
}

func Start(proxy Proxy) error {
	if servicePlugin == nil {
		return fmt.Errorf("no service plugin registered")
	}

	return servicePlugin.Start(proxy)
}
