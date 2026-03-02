package prompt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/ntthienan0507-web/create-go-api/internal/scaffold"
)

// Run displays interactive prompts and returns the user's choices.
func Run() (*scaffold.ProjectConfig, error) {
	cfg := &scaffold.ProjectConfig{
		ServerPort:          8080,
		IncludeSampleModule: true,
	}

	var portStr string = "8080"

	form := huh.NewForm(
		// Group 1: Project identity
		huh.NewGroup(
			huh.NewInput().
				Title("Go module path").
				Description("e.g. github.com/myorg/my-api").
				Value(&cfg.ModulePath).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("module path is required")
					}
					if !strings.Contains(s, "/") {
						return errors.New("module path must contain at least one /")
					}
					if strings.Contains(s, " ") {
						return errors.New("module path must not contain spaces")
					}
					return nil
				}),

			huh.NewInput().
				Title("Server port").
				Description("Default: 8080").
				Value(&portStr).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					n, err := strconv.Atoi(s)
					if err != nil {
						return errors.New("port must be a number")
					}
					if n < 1 || n > 65535 {
						return errors.New("port must be between 1 and 65535")
					}
					return nil
				}),
		),

		// Group 2: Database driver
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Database driver").
				Options(
					huh.NewOption("SQLC (type-safe SQL)", "sqlc"),
					huh.NewOption("GORM (ORM)", "gorm"),
					huh.NewOption("Both (switchable via config)", "both"),
				).
				Value(&cfg.DBDriver),
		),

		// Group 3: Auth provider
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Auth provider").
				Options(
					huh.NewOption("JWT (local HMAC tokens)", "jwt"),
					huh.NewOption("Keycloak (OIDC)", "keycloak"),
					huh.NewOption("Both (switchable via config)", "both"),
				).
				Value(&cfg.AuthProvider),
		),

		// Group 4: Sample module
		huh.NewGroup(
			huh.NewConfirm().
				Title("Include sample User module?").
				Description("Generates a complete CRUD module as reference").
				Value(&cfg.IncludeSampleModule),
		),
	)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("prompt cancelled: %w", err)
	}

	// Parse port
	if portStr != "" {
		n, _ := strconv.Atoi(portStr)
		if n > 0 {
			cfg.ServerPort = n
		}
	}

	cfg.Derive()
	return cfg, nil
}
