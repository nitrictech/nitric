package api

import (
	"encoding/json"
	"io"
)

type AuthDetails struct {
	WorkOS WorkOSDetails `json:"workos"`
}

type WorkOSDetails struct {
	ClientID    string `json:"client_id"`
	ApiHostname string `json:"api_hostname"`
}

func (c *NitricApiClient) GetWorkOSPublicDetails() (*WorkOSDetails, error) {
	response, err := c.get("/auth/details")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var authDetails AuthDetails
	err = json.Unmarshal(body, &authDetails)
	if err != nil {
		return nil, err
	}

	return &authDetails.WorkOS, nil
}
