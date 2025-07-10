package app

import (
	"errors"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/workos"
	"github.com/samber/do/v2"
)

type AuthApp struct {
	auth *workos.WorkOSAuth
}

func NewAuthApp(injector do.Injector) (*AuthApp, error) {
	auth := do.MustInvoke[*workos.WorkOSAuth](injector)
	return &AuthApp{auth: auth}, nil
}

// Login handles the login command logic
func (c *AuthApp) Login() {
	fmt.Printf("\n%s Logging in...\n", style.Purple(icons.Lightning+" Nitric"))

	user, err := c.auth.Login()
	if err != nil {
		fmt.Printf("\n%s Error logging in: %s\n", style.Red(icons.Cross), err)
		return
	}

	fmt.Printf("\n%s Logged in as %s\n", style.Green(icons.Check), style.Teal(user.FirstName))
}

// Logout handles the logout command logic
func (c *AuthApp) Logout() {
	fmt.Printf("\n%s Logging out...\n", style.Purple(icons.Lightning+" Nitric"))

	err := c.auth.Logout()
	if err != nil {
		if !errors.Is(err, workos.ErrNotFound) {
			fmt.Printf("\n%s Error logging out: %s\n", style.Red(icons.Cross), err)
			return
		}
	}

	fmt.Printf("\n%s Logged out successfully\n", style.Green(icons.Check))
}

// AccessToken handles the access token command logic
func (c *AuthApp) AccessToken() {
	token, err := c.auth.GetAccessToken()
	if err != nil {
		fmt.Printf("\n%s Error getting access token: %s\n", style.Red(icons.Cross), err)
		return
	}

	fmt.Printf("\n%s Access token: %s\n", style.Green(icons.Check), token)
}
