package app

import (
	"errors"
	"fmt"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/nitrictech/nitric/cli/pkg/tui"
	"github.com/nitrictech/nitric/cli/pkg/tui/ask"
	"github.com/samber/do/v2"
)

type TeamApp struct {
	apiClient *api.NitricApiClient
	auth      *workos.WorkOSAuth
	styles    tui.AppStyles
}

func NewTeamApp(injector do.Injector) (*TeamApp, error) {
	apiClient, err := api.NewNitricApiClient(injector)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	auth := do.MustInvoke[*workos.WorkOSAuth](injector)

	styles := tui.NewAppStyles()

	return &TeamApp{
		apiClient: apiClient,
		auth:      auth,
		styles:    styles,
	}, nil
}

func (t *TeamApp) SwitchTeam(teamSlug string) error {
	allTeams, err := t.apiClient.GetUserTeams()
	if err != nil {
		if errors.Is(err, api.ErrUnauthenticated) {
			fmt.Println("Please login first, using the", t.styles.Emphasize.Render(version.GetCommand("login")), "command")
			return nil
		}
		fmt.Printf("Failed to get teams: %v\n", err)
		return nil
	}

	if len(allTeams) == 0 {
		fmt.Println("No teams found. Create a team first to continue.")
		return nil
	}

	if teamSlug != "" {
		return t.switchToTeamBySlug(allTeams, teamSlug)
	}

	return t.showInteractiveTeamPicker(allTeams)
}

func (t *TeamApp) switchToTeamBySlug(teams []api.Team, slug string) error {
	var targetTeam *api.Team
	for _, team := range teams {
		if team.Slug == slug {
			targetTeam = &team
			break
		}
	}

	if targetTeam == nil {
		return fmt.Errorf("team not found: %s", slug)
	}

	if targetTeam.IsCurrent {
		fmt.Printf("%s Already using team: %s\n", style.Gray(icons.Check), style.Teal(targetTeam.Name))
		return nil
	}

	return t.performTeamSwitch(targetTeam)
}

func (t *TeamApp) showInteractiveTeamPicker(teams []api.Team) error {
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})

	var currentTeam *api.Team
	for i := range teams {
		if teams[i].IsCurrent {
			currentTeam = &teams[i]
			break
		}
	}

	if currentTeam != nil {
		currentStyle := lipgloss.NewStyle().
			Foreground(colors.Teal).
			Bold(true)
		fmt.Printf("Current team: %s\n\n", currentStyle.Render(currentTeam.Name))
	}

	teamMap := make(map[string]*api.Team)
	optionLabels := make([]string, 0, len(teams))

	for i := range teams {
		team := &teams[i]
		label := team.Name
		if team.IsCurrent {
			label = fmt.Sprintf("%s (current)", team.Name)
		}
		teamMap[label] = team
		optionLabels = append(optionLabels, label)
	}

	var selectedLabel string
	err := ask.NewSelect[string]().
		Title("Select a team:").
		Options(huh.NewOptions(optionLabels...)...).
		Value(&selectedLabel).
		Height(len(teams) + 2). // +2 for extra spacing at the bottom
		Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		return fmt.Errorf("failed to get team selection: %w", err)
	}

	selectedTeam := teamMap[selectedLabel]
	if selectedTeam.IsCurrent {
		fmt.Printf("%s Already using team: %s\n", style.Gray(icons.Check), style.Teal(selectedTeam.Name))
		return nil
	}

	return t.performTeamSwitch(selectedTeam)
}

func (t *TeamApp) performTeamSwitch(team *api.Team) error {
	fmt.Printf("Switching to team: %s\n", style.Teal(team.Name))

	err := t.auth.RefreshTokenForOrganization(team.WorkOsID)
	if err != nil {
		fmt.Printf("%s Failed to refresh token for organization: %v\n", style.Red(icons.Cross), err)
		fmt.Printf("Try running %s to re-authenticate\n", style.Teal(version.CommandName+" login"))
		return nil
	}

	fmt.Printf("%s Successfully switched to team: %s\n", style.Green(icons.Check), style.Teal(team.Name))
	return nil
}
