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

func (c *NitricApiClient) GetUserTeams() ([]Team, error) {
	response, err := c.get("/me/teams", true)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user teams: %s", response.Status)
	}

	var teamsResponse GetTeamsResponse
	if err := json.NewDecoder(response.Body).Decode(&teamsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teams: %v", err)
	}

	return teamsResponse.Teams, nil
}
