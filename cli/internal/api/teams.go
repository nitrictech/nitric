package api

import (
	"encoding/json"
	"fmt"
)

type Team struct {
	Slug      string `json:"slug"`
	WorkOsID  string `json:"workOsId"`
	ImageUrl  string `json:"imageUrl"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	IsCurrent bool   `json:"isCurrent"`
}

type GetTeamsResponse struct {
	Teams []Team `json:"teams"`
}

func (c *SugaApiClient) GetUserTeams() ([]Team, error) {
	response, err := c.get("/me/teams", true)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == 401 || response.StatusCode == 403 {
		return nil, ErrUnauthenticated
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user teams: %s", response.Status)
	}

	var teamsResponse GetTeamsResponse
	if err := json.NewDecoder(response.Body).Decode(&teamsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teams: %v", err)
	}

	return teamsResponse.Teams, nil
}

type SwitchTeamRequest struct {
	OrganizationID string `json:"organizationId"`
}

type SwitchTeamResponse struct {
	RedirectURL string `json:"redirectUrl"`
}

func (c *SugaApiClient) SwitchTeam(organizationId string) (*SwitchTeamResponse, error) {
	requestBody := SwitchTeamRequest{
		OrganizationID: organizationId,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	response, err := c.post("/me/teams/switch", true, jsonData)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == 401 || response.StatusCode == 403 {
		return nil, ErrUnauthenticated
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to switch team: %s", response.Status)
	}

	var switchResponse SwitchTeamResponse
	if err := json.NewDecoder(response.Body).Decode(&switchResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal switch team response: %v", err)
	}

	return &switchResponse, nil
}
