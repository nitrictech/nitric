package terraform

import (
	"fmt"
	"io"
	"maps"
	"os"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type PlatformSpec struct {
	Name            string                       `json:"name" yaml:"name"`
	ServicesSpec    NitricServiceSpec            `json:"services" yaml:"services"`
	BucketsSpec     NitricResourceSpec           `json:"buckets,omitempty" yaml:"buckets,omitempty"`
	TopicsSpec      NitricResourceSpec           `json:"topics,omitempty" yaml:"topics,omitempty"`
	EntrypointsSpec NitricResourceSpec           `json:"entrypoints" yaml:"entrypoints"`
	Infra           map[string]InfraResourceSpec `json:"infra" yaml:"infra"`
}

func (p PlatformSpec) GetServiceSpec(subtype string) (ServiceSpec, error) {
	spec := &p.ServicesSpec

	if subtype != "" {
		subspec, ok := spec.Subtypes[subtype]
		if !ok {
			return ServiceSpec{}, fmt.Errorf("platform %s does not define subtype %s for %s, available subtypes: %v", p.Name, subtype, typ, maps.Keys(spec.Subtypes))
		}

		return subspec, nil
	}

	return spec.ServiceSpec, nil
}

func (p PlatformSpec) GetResourceSpecForTypes(typ string, subtype string) (ResourceSpec, error) {
	var spec *NitricResourceSpec
	switch typ {

	case "entrypoint":
		spec = &p.EntrypointsSpec
	case "bucket":
		spec = &p.BucketsSpec
	case "topic":
		spec = &p.TopicsSpec
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

func PlatformSpecFromReader(reader io.Reader) (*PlatformSpec, error) {
	var spec PlatformSpec

	byt, err := afero.ReadAll(reader)
	if err != nil {
		return &PlatformSpec{}, nil
	}

	err = yaml.Unmarshal(byt, &spec)

	return &spec, err
}

func PlatformSpecFromFile(fs afero.Fs, filePath string) (*PlatformSpec, error) {
	file, err := fs.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return &PlatformSpec{}, fmt.Errorf("failed to read platform spec file %s: %w", filePath, err)
	}

	return PlatformSpecFromReader(file)
}

type PlatformReferencePrefix string

const (
	PlatformReferencePrefix_File  = "file:"
	PlatformReferencePrefix_Https = "https://"
	PlatformReferencePrefix_Git   = "git+"
)

func PlatformFromId(fs afero.Fs, platformId string, repositories ...PlatformRepository) (*PlatformSpec, error) {
	if strings.HasPrefix(platformId, PlatformReferencePrefix_File) {
		return PlatformSpecFromFile(fs, strings.Replace(platformId, PlatformReferencePrefix_File, "", -1))
	} else if strings.HasPrefix(platformId, PlatformReferencePrefix_Https) || strings.HasPrefix(platformId, PlatformReferencePrefix_Git) {
		return nil, fmt.Errorf("platform %s is not supported yet", platformId)
	}

	for _, repository := range repositories {
		platform, err := repository.GetPlatform(platformId)
		if err != nil {
			continue
		}

		return platform, nil
	}

	// TODO: check for close matches and list available platforms
	return nil, fmt.Errorf("platform %s not found in any repository", platformId)
}

type ResourceSpec struct {
	PluginId   string                 `json:"plugin" yaml:"plugin"`
	Properties map[string]interface{} `json:"properties" yaml:"properties"`
}

type ServiceSpec struct {
	ResourceSpec `json:",inline" yaml:",inline"`
	Identities   map[string]InfraResourceSpec `json:"identities" yaml:"identities"`
}

type NitricServiceSpec struct {
	ServiceSpec `json:",inline" yaml:",inline"`
	Subtypes    map[string]ServiceSpec `json:"subtypes" yaml:"subtypes"`
}
type NitricResourceSpec struct {
	ResourceSpec `json:",inline" yaml:",inline"`
	Subtypes     map[string]ResourceSpec `json:"subtypes" yaml:"subtypes"`
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
	ResourceSpec `json:",inline" yaml:",inline"`
}
