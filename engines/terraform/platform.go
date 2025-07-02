package terraform

import (
	"fmt"
	"io"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type PlatformSpec struct {
	Name string `json:"name" yaml:"name"`

	Libraries map[string]string `json:"libraries" yaml:"libraries"`

	Variables map[string]Variable `json:"variables" yaml:"variables,omitempty"`

	ServiceBlueprints    map[string]*ServiceBlueprint  `json:"services" yaml:"services"`
	BucketBlueprints     map[string]*ResourceBlueprint `json:"buckets,omitempty" yaml:"buckets,omitempty"`
	TopicBlueprints      map[string]*ResourceBlueprint `json:"topics,omitempty" yaml:"topics,omitempty"`
	DatabaseBlueprints   map[string]*ResourceBlueprint `json:"databases,omitempty" yaml:"databases,omitempty"`
	EntrypointBlueprints map[string]*ResourceBlueprint `json:"entrypoints" yaml:"entrypoints"`
	InfraSpecs           map[string]*ResourceBlueprint `json:"infra" yaml:"infra"`
}

type Variable struct {
	Type        string
	Description string
	Default     interface{}
}

type Library struct {
	Team    string `json:"team" yaml:"team"`
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

func (p PlatformSpec) GetLibrary(name string) (*Library, error) {
	library, ok := p.Libraries[name]
	if !ok {
		return nil, fmt.Errorf("library %s not found in platform spec", name)
	}

	pattern := `^(?P<team>[^/]+)/(?P<library>[^@]+)@(?P<version>.+)$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(library)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid package format: %s", library)
	}

	team := matches[re.SubexpIndex("team")]
	libName := matches[re.SubexpIndex("library")]
	version := matches[re.SubexpIndex("version")]

	return &Library{Team: team, Name: libName, Version: version}, nil
}

func (p PlatformSpec) GetServiceBlueprint(intentSubType string) (*ServiceBlueprint, error) {
	spec := p.ServiceBlueprints

	if intentSubType == "" {
		intentSubType = "default"
	}

	concreteSpec, ok := spec[intentSubType]
	if !ok {
		return nil, fmt.Errorf("platform %s does not define a %s type for services, available types: %v", p.Name, intentSubType, slices.Collect(maps.Keys(spec)))
	}

	return concreteSpec, nil
}

func (p PlatformSpec) GetResourceBlueprintsForType(typ string) (map[string]*ResourceBlueprint, error) {
	var specs map[string]*ResourceBlueprint

	switch typ {
	case "service":
		specs = map[string]*ResourceBlueprint{}
		for name, blueprint := range p.ServiceBlueprints {
			specs[name] = blueprint.ResourceBlueprint
		}
	case "entrypoint":
		specs = p.EntrypointBlueprints
	case "bucket":
		specs = p.BucketBlueprints
	case "topic":
		specs = p.TopicBlueprints
	default:
		return nil, fmt.Errorf("failed to resolve resource blueprint, no type %s in platform spec", typ)
	}

	return specs, nil
}

func (p PlatformSpec) GetResourceBlueprint(intentType string, intentSubType string) (*ResourceBlueprint, error) {
	if intentSubType == "" {
		intentSubType = "default"
	}

	var spec *ResourceBlueprint
	switch intentType {
	case "service":
		if serviceBlueprint, ok := p.ServiceBlueprints[intentSubType]; ok {
			spec = serviceBlueprint.ResourceBlueprint
		}
	case "entrypoint":
		spec, _ = p.EntrypointBlueprints[intentSubType]
	case "bucket":
		spec, _ = p.BucketBlueprints[intentSubType]
	case "topic":
		spec, _ = p.TopicBlueprints[intentSubType]
	case "database":
		spec, _ = p.DatabaseBlueprints[intentSubType]
	default:
		return nil, fmt.Errorf("failed to resolve resource blueprint, no type %s known in platform spec", intentType)
	}

	if spec == nil {
		return nil, fmt.Errorf("platform %s does not define a '%s' %s type", p.Name, intentSubType, intentType)
	}

	return spec, nil
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

type ResourceBlueprint struct {
	PluginId   string                 `json:"plugin" yaml:"plugin"`
	Properties map[string]interface{} `json:"properties" yaml:"properties"`
	DependsOn  []string               `json:"depends_on" yaml:"depends_on,omitempty"`
	Variables  map[string]Variable    `json:"variables" yaml:"variables,omitempty"`
}

type IdentitiesBlueprint struct {
	Identities []ResourceBlueprint `json:"identities" yaml:"identities"`
}

func (i IdentitiesBlueprint) GetIdentities() []ResourceBlueprint {
	if i.Identities == nil {
		return []ResourceBlueprint{}
	}
	return i.Identities
}

type Identifiable interface {
	GetIdentity(string) (*ResourceBlueprint, error)
	GetIdentities() map[string]ResourceBlueprint
}

type ServiceBlueprint struct {
	*ResourceBlueprint   `json:",inline" yaml:",inline"`
	*IdentitiesBlueprint `json:",inline" yaml:",inline"`
}
