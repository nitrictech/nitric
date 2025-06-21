package api

import (
	"encoding/json"
	"io"
)

type WorkOSPublicDetails struct {
	ClientID    string `json:"client_id"`
	ApiHostname string `json:"api_hostname"`
}

func (c *NitricApiClient) GetWorkOSPublicDetails() (*WorkOSPublicDetails, error) {
	response, err := c.get("/auth/public/workos")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var workOSPublicDetails WorkOSPublicDetails
	err = json.Unmarshal(body, &workOSPublicDetails)
	if err != nil {
		return nil, err
	}

	return &workOSPublicDetails, nil
}
