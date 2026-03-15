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
	var infraChoices, featureChoices, serviceChoices []string

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

		// Group 5: Swagger
		huh.NewGroup(
			huh.NewConfirm().
				Title("Include Swagger/OpenAPI docs?").
				Description("Adds swag annotations and /swagger/index.html endpoint").
				Value(&cfg.IncludeSwagger),
		),

		// Group 6: Infrastructure
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Infrastructure (space to select)").
				Options(
					huh.NewOption("Redis — cache, sessions, pub/sub", "redis"),
					huh.NewOption("Kafka — event streaming, outbox pattern", "kafka"),
					huh.NewOption("Encryption — AES-256 for PII fields", "encryption"),
				).
				Value(&infraChoices),
		),

		// Group 7: Extra features
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Features (space to select)").
				Options(
					huh.NewOption("Cron — scheduled background jobs", "cron"),
					huh.NewOption("WebSocket — real-time communication", "ws"),
					huh.NewOption("OpenTelemetry — distributed tracing", "otel"),
				).
				Value(&featureChoices),
		),

		// Group 8: External services
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("External services (space to select)").
				Options(
					huh.NewOption("SendGrid — transactional email", "sendgrid"),
					huh.NewOption("Stripe — payment processing", "stripe"),
					huh.NewOption("IceWarp — mail server (XML API)", "icewarp"),
					huh.NewOption("Firebase — push notifications", "firebase"),
					huh.NewOption("Elasticsearch — search & analytics", "elasticsearch"),
				).
				Value(&serviceChoices),
		),

		// Group 9: DevOps
		huh.NewGroup(
			huh.NewConfirm().Title("Include Kubernetes manifests?").Value(&cfg.IncludeK8s),
			huh.NewConfirm().Title("Include SonarQube config?").Value(&cfg.IncludeSonar),
			huh.NewSelect[string]().
				Title("CI/CD provider").
				Options(
					huh.NewOption("None", ""),
					huh.NewOption("GitHub Actions", "github"),
					huh.NewOption("GitLab CI", "gitlab"),
					huh.NewOption("Both", "both"),
				).
				Value(&cfg.CIProvider),
		),
	)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("prompt cancelled: %w", err)
	}

	// Map multi-select choices to config
	for _, c := range infraChoices {
		switch c {
		case "redis":      cfg.IncludeRedis = true
		case "kafka":      cfg.IncludeKafka = true
		case "encryption": cfg.IncludeEncryption = true
		}
	}
	for _, c := range featureChoices {
		switch c {
		case "cron": cfg.IncludeCron = true
		case "ws":   cfg.IncludeWebSocket = true
		case "otel": cfg.IncludeOTEL = true
		}
	}
	for _, c := range serviceChoices {
		switch c {
		case "sendgrid":      cfg.IncludeSendGrid = true
		case "stripe":        cfg.IncludeStripe = true
		case "icewarp":       cfg.IncludeIceWarp = true
		case "firebase":      cfg.IncludeFirebase = true
		case "elasticsearch": cfg.IncludeElasticsearch = true
		}
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
