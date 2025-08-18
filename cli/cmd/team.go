package cmd

import (
	"github.com/nitrictech/nitric/cli/pkg/app"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// NewTeamCmd creates the team command
func NewTeamCmd(injector do.Injector) *cobra.Command {
	teamCmd := &cobra.Command{
		Use:   "team [team-slug]",
		Short: "Switch between teams",
		Long: `Switch between teams or list available teams.

When run without arguments, displays an interactive list of available teams.
When run with a team slug, switches directly to that team.`,
		Example: `
# Show interactive team switcher
nitric team

# Switch to team by slug
nitric team my-team-slug
		`,
		Run: func(cmd *cobra.Command, args []string) {
			teamApp, err := app.NewTeamApp(injector)
			if err != nil {
				cobra.CheckErr(err)
			}

			teamSlug := ""
			if len(args) == 1 {
				teamSlug = args[0]
			}

			err = teamApp.SwitchTeam(teamSlug)
			cobra.CheckErr(err)
		},
	}

	return teamCmd
}
