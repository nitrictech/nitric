package terraform

import (
	"fmt"
	"maps"
)

type PlatformSpec struct {
	Name            string                       `json:"name"`
	ServicesSpec    NitricResourceSpec           `json:"services"`
	EntrypointsSpec NitricResourceSpec           `json:"entrypoints"`
	Infra           map[string]InfraResourceSpec `json:"infra"`
}

func (p PlatformSpec) GetResourceSpecForTypes(typ string, subtype string) (ResourceSpec, error) {
	var spec *NitricResourceSpec
	switch typ {
	case "service":
		spec = &p.ServicesSpec
	case "entrypoint":
		spec = &p.EntrypointsSpec
	default:
		return ResourceSpec{}, fmt.Errorf("no type %s known in platform spec", typ)
	}

	if subtype != "" {
		subspec, ok := spec.Subtypes[subtype]
		if !ok {
			return ResourceSpec{}, fmt.Errorf("platform %s does not define subtype %s for %s, available subtypes: %v", p.Name, subtype, typ, maps.Keys(spec.Subtypes))
		}

		return subspec, nil
	}

	return spec.ResourceSpec, nil
}

type ResourceSpec struct {
	PluginId   string                 `json:"plugin"`
	Properties map[string]interface{} `json:"properties"`
}

type NitricResourceSpec struct {
	ResourceSpec `json:",inline"`
	Subtypes     map[string]ResourceSpec `json:"subtypes"`
}

func (r NitricResourceSpec) GetResourceProperties(subtype string) map[string]interface{} {
	if subtype != "" {
		return r.Subtypes[subtype].Properties
	}

	return r.Properties
}

func (r NitricResourceSpec) GetPlugin(subtype string) (string, error) {
	if subtype != "" {
		if _, ok := r.Subtypes[subtype]; !ok {
			return "", fmt.Errorf("subtype %s not found", subtype)
		}
		return r.Subtypes[subtype].PluginId, nil
	}
	return r.PluginId, nil
}

type InfraResourceSpec struct {
	ResourceSpec `json:",inline"`
}
