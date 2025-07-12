package platforms

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/engines/terraform"
)

type PlatformRepository struct {
	apiClient *api.NitricApiClient
}

var _ terraform.PlatformRepository = (*PlatformRepository)(nil)

func (r *PlatformRepository) GetPlatform(name string) (*terraform.PlatformSpec, error) {
	// Split the name into team, lib, and revision using a regex <team>/<lib>@<revision>
	re := regexp.MustCompile(`^(?P<team>[^/]+)/(?P<platform>[^@]+)@(?P<revision>\d+)$`)
	matches := re.FindStringSubmatch(name)

	if matches == nil {
		return nil, fmt.Errorf("invalid platform name format: %s. Expected format: <team>/<lib>@<revision> e.g. nitric/aws@1", name)
	}

	// Extract named groups
	team := matches[re.SubexpIndex("team")]
	platform := matches[re.SubexpIndex("platform")]
	revisionStr := matches[re.SubexpIndex("revision")]

	// Convert revision string to integer
	revision, err := strconv.Atoi(revisionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid revision format: %s. Expected integer", revisionStr)
	}

	platformSpec, err := r.apiClient.GetPlatform(team, platform, revision)
	// If its a 404, then return platform not found error
	if err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return nil, terraform.ErrPlatformNotFound
		} else if errors.Is(err, api.ErrUnauthenticated) {
			return nil, terraform.ErrUnauthenticated
		}

		// return the original error to the engine
		return nil, err
	}

	return platformSpec, nil
}

func NewPlatformRepository(apiClient *api.NitricApiClient) *PlatformRepository {
	return &PlatformRepository{
		apiClient: apiClient,
	}
}
