package core

import (
	"io"
	"net/http"

	"github.com/nitrictech/nitric/pkg/utils"
)

type GcpProvider interface {
	// GetServiceAccountEmail for google cloud projects
	GetServiceAccountEmail() (string, error)
	GetProjectID() (string, error)
}

type gcpProviderImpl struct {
	serviceAccountEmail string
	projectID           string
}

const (
	metadataFlavorKey      = "Metadata-Flavor"
	metadataFlavorValue    = "Google"
	serviceAccountEnv      = "SERVICE_ACCOUNT_EMAIL"
	serviceAccountEmailUri = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/email"
	projectIdEnv           = "GOOGLE_PROJECT_ID"
	projectIdUri           = "http://metadata.google.internal/computeMetadata/v1/project/project-id"
)

func createMetadataRequest(uri string) (*http.Request, error) {
	req, err := http.NewRequest("GET", projectIdUri, nil)
	if err != nil {
		return nil, err
	}

	// Add the correct header
	req.Header.Add(metadataFlavorKey, metadataFlavorValue)

	return req, nil
}

func (g *gcpProviderImpl) GetProjectID() (string, error) {
	if g.projectID == "" {
		if env := utils.GetEnv(projectIdEnv, ""); env != "" {
			return env, nil
		}

		req, err := createMetadataRequest(projectIdUri)
		if err != nil {
			return "", err
		}

		// read the response as the service account email
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		projectIdBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		g.projectID = string(projectIdBytes)
	}

	return g.projectID, nil
}

func (g *gcpProviderImpl) GetServiceAccountEmail() (string, error) {
	if g.serviceAccountEmail == "" {
		if env := utils.GetEnv(g.serviceAccountEmail, ""); env != "" {
			return env, nil
		}

		req, err := createMetadataRequest(serviceAccountEmailUri)
		if err != nil {
			return "", err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		// read the response as the service account email
		emailBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		g.serviceAccountEmail = string(emailBytes)
	}

	return g.serviceAccountEmail, nil
}

func New() (GcpProvider, error) {
	return &gcpProviderImpl{}, nil
}
