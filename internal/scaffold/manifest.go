package scaffold

// FileRule defines when a template file should be included.
type FileRule struct {
	TmplPath    string                        // path inside templates/ FS
	OutPath     string                        // static output path
	OutPathFunc func(*ProjectConfig) string   // dynamic output path (overrides OutPath)
	Include     func(*ProjectConfig) bool     // nil = always include
}

// manifest lists every template file and its inclusion predicate.
var manifest = []FileRule{
	// --- Always included (core skeleton) ---
	{TmplPath: "cmd/app/main.go.tmpl",
		OutPathFunc: func(c *ProjectConfig) string { return "cmd/" + c.ProjectName + "/main.go" }},
	{TmplPath: "internal/config/config.go.tmpl", OutPath: "internal/config/config.go"},
	{TmplPath: "internal/app/app.go.tmpl", OutPath: "internal/app/app.go"},
	{TmplPath: "internal/middleware/auth.go.tmpl", OutPath: "internal/middleware/auth.go"},
	{TmplPath: "internal/middleware/cors.go.tmpl", OutPath: "internal/middleware/cors.go"},
	{TmplPath: "internal/middleware/logger.go.tmpl", OutPath: "internal/middleware/logger.go"},
	{TmplPath: "internal/middleware/recovery.go.tmpl", OutPath: "internal/middleware/recovery.go"},
	{TmplPath: "internal/response/response.go.tmpl", OutPath: "internal/response/response.go"},
	{TmplPath: "internal/response/pagination.go.tmpl", OutPath: "internal/response/pagination.go"},
	{TmplPath: "internal/auth/provider.go.tmpl", OutPath: "internal/auth/provider.go"},
	{TmplPath: "internal/auth/factory.go.tmpl", OutPath: "internal/auth/factory.go"},
	{TmplPath: "internal/database/migration.go.tmpl", OutPath: "internal/database/migration.go"},
	{TmplPath: "db/migrations/00001_create_users.sql.tmpl", OutPath: "db/migrations/00001_create_users.sql"},
	{TmplPath: "pkg/logger/logger.go.tmpl", OutPath: "pkg/logger/logger.go"},
	{TmplPath: "Makefile.tmpl", OutPath: "Makefile"},
	{TmplPath: "Dockerfile.tmpl", OutPath: "Dockerfile"},
	{TmplPath: "dot-env.example.tmpl", OutPath: ".env.example"},
	{TmplPath: "dot-gitignore.tmpl", OutPath: ".gitignore"},
	{TmplPath: "go.mod.tmpl", OutPath: "go.mod"},

	// --- Swagger ---
	{TmplPath: "docs/dot-gitkeep.tmpl", OutPath: "docs/.gitkeep",
		Include: func(c *ProjectConfig) bool { return c.IncludeSwagger }},

	// --- SQLC-only files ---
	{TmplPath: "internal/database/postgres.go.tmpl", OutPath: "internal/database/postgres.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSQLC }},
	{TmplPath: "internal/database/store.go.tmpl", OutPath: "internal/database/store.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSQLC }},
	{TmplPath: "db/queries/users.sql.tmpl", OutPath: "db/queries/users.sql",
		Include: func(c *ProjectConfig) bool { return c.IncludeSQLC && c.IncludeSampleModule }},
	{TmplPath: "db/sqlc/db.go.tmpl", OutPath: "db/sqlc/db.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSQLC }},
	{TmplPath: "db/sqlc/models.go.tmpl", OutPath: "db/sqlc/models.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSQLC && c.IncludeSampleModule }},
	{TmplPath: "db/sqlc/querier.go.tmpl", OutPath: "db/sqlc/querier.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSQLC && c.IncludeSampleModule }},
	{TmplPath: "db/sqlc/users.sql.go.tmpl", OutPath: "db/sqlc/users.sql.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSQLC && c.IncludeSampleModule }},
	{TmplPath: "sqlc.yaml.tmpl", OutPath: "sqlc.yaml",
		Include: func(c *ProjectConfig) bool { return c.IncludeSQLC }},

	// --- GORM-only files ---
	{TmplPath: "internal/database/gorm.go.tmpl", OutPath: "internal/database/gorm.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeGORM }},

	// --- Auth conditional ---
	{TmplPath: "internal/auth/jwt.go.tmpl", OutPath: "internal/auth/jwt.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeJWT }},
	{TmplPath: "internal/auth/keycloak.go.tmpl", OutPath: "internal/auth/keycloak.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKeycloak }},

	// --- Sample module (user) ---
	{TmplPath: "internal/module/user/types.go.tmpl", OutPath: "internal/module/user/types.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSampleModule }},
	{TmplPath: "internal/module/user/repository.go.tmpl", OutPath: "internal/module/user/repository.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSampleModule }},
	{TmplPath: "internal/module/user/service.go.tmpl", OutPath: "internal/module/user/service.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSampleModule }},
	{TmplPath: "internal/module/user/handler.go.tmpl", OutPath: "internal/module/user/handler.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSampleModule }},
	{TmplPath: "internal/module/user/routes.go.tmpl", OutPath: "internal/module/user/routes.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSampleModule }},
	{TmplPath: "internal/module/user/repository_sqlc.go.tmpl", OutPath: "internal/module/user/repository_sqlc.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSampleModule && c.IncludeSQLC }},
	{TmplPath: "internal/module/user/repository_gorm.go.tmpl", OutPath: "internal/module/user/repository_gorm.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSampleModule && c.IncludeGORM }},

	// --- Redis ---
	{TmplPath: "pkg/database/redis.go.tmpl", OutPath: "pkg/database/redis.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeRedis }},
	{TmplPath: "pkg/cache/redis.go.tmpl", OutPath: "pkg/cache/redis.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeRedis }},

	// --- Kafka + Outbox ---
	{TmplPath: "pkg/broker/broker.go.tmpl", OutPath: "pkg/broker/broker.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/topic.go.tmpl", OutPath: "pkg/broker/topic.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/config.go.tmpl", OutPath: "pkg/broker/config.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/producer.go.tmpl", OutPath: "pkg/broker/producer.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/consumer.go.tmpl", OutPath: "pkg/broker/consumer.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/dispatcher.go.tmpl", OutPath: "pkg/broker/dispatcher.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/batcher.go.tmpl", OutPath: "pkg/broker/batcher.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/outbox.go.tmpl", OutPath: "pkg/broker/outbox.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/relay.go.tmpl", OutPath: "pkg/broker/relay.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/idempotency.go.tmpl", OutPath: "pkg/broker/idempotency.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/errors.go.tmpl", OutPath: "pkg/broker/errors.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "pkg/broker/scram.go.tmpl", OutPath: "pkg/broker/scram.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "db/migrations/00002_create_outbox.sql.tmpl", OutPath: "db/migrations/00002_create_outbox.sql",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},
	{TmplPath: "db/migrations/00003_create_processed_events.sql.tmpl", OutPath: "db/migrations/00003_create_processed_events.sql",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka }},

	// --- Kafka sample events (only with sample module) ---
	{TmplPath: "modules/user/events.go.tmpl", OutPath: "modules/user/events.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeKafka && c.IncludeSampleModule }},

	// --- Encryption ---
	{TmplPath: "pkg/crypto/crypto.go.tmpl", OutPath: "pkg/crypto/crypto.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeEncryption }},

	// --- Cron ---
	{TmplPath: "pkg/cron/cron.go.tmpl", OutPath: "pkg/cron/cron.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeCron }},
	{TmplPath: "cmd/cron.go.tmpl", OutPath: "cmd/cron.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeCron }},

	// --- WebSocket ---
	{TmplPath: "pkg/ws/message.go.tmpl", OutPath: "pkg/ws/message.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeWebSocket }},
	{TmplPath: "pkg/ws/client.go.tmpl", OutPath: "pkg/ws/client.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeWebSocket }},
	{TmplPath: "pkg/ws/hub.go.tmpl", OutPath: "pkg/ws/hub.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeWebSocket }},
	{TmplPath: "pkg/ws/handler.go.tmpl", OutPath: "pkg/ws/handler.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeWebSocket }},

	// --- OpenTelemetry ---
	{TmplPath: "pkg/tracing/tracing.go.tmpl", OutPath: "pkg/tracing/tracing.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeOTEL }},

	// --- External: SendGrid ---
	{TmplPath: "pkg/external/sendgrid/client.go.tmpl", OutPath: "pkg/external/sendgrid/client.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSendGrid }},
	{TmplPath: "pkg/external/sendgrid/types.go.tmpl", OutPath: "pkg/external/sendgrid/types.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSendGrid }},
	{TmplPath: "pkg/external/sendgrid/errors.go.tmpl", OutPath: "pkg/external/sendgrid/errors.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeSendGrid }},

	// --- External: Stripe ---
	{TmplPath: "pkg/external/stripe/client.go.tmpl", OutPath: "pkg/external/stripe/client.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeStripe }},
	{TmplPath: "pkg/external/stripe/types.go.tmpl", OutPath: "pkg/external/stripe/types.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeStripe }},
	{TmplPath: "pkg/external/stripe/errors.go.tmpl", OutPath: "pkg/external/stripe/errors.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeStripe }},

	// --- External: IceWarp ---
	{TmplPath: "pkg/external/icewarp/client.go.tmpl", OutPath: "pkg/external/icewarp/client.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeIceWarp }},
	{TmplPath: "pkg/external/icewarp/types.go.tmpl", OutPath: "pkg/external/icewarp/types.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeIceWarp }},
	{TmplPath: "pkg/external/icewarp/errors.go.tmpl", OutPath: "pkg/external/icewarp/errors.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeIceWarp }},

	// --- External: Firebase ---
	{TmplPath: "pkg/external/firebase/client.go.tmpl", OutPath: "pkg/external/firebase/client.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeFirebase }},
	{TmplPath: "pkg/external/firebase/types.go.tmpl", OutPath: "pkg/external/firebase/types.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeFirebase }},
	{TmplPath: "pkg/external/firebase/errors.go.tmpl", OutPath: "pkg/external/firebase/errors.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeFirebase }},

	// --- External: Elasticsearch ---
	{TmplPath: "pkg/external/elasticsearch/client.go.tmpl", OutPath: "pkg/external/elasticsearch/client.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeElasticsearch }},
	{TmplPath: "pkg/external/elasticsearch/types.go.tmpl", OutPath: "pkg/external/elasticsearch/types.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeElasticsearch }},
	{TmplPath: "pkg/external/elasticsearch/errors.go.tmpl", OutPath: "pkg/external/elasticsearch/errors.go",
		Include: func(c *ProjectConfig) bool { return c.IncludeElasticsearch }},

	// --- Kubernetes ---
	{TmplPath: "deployments/k8s/base/kustomization.yaml.tmpl", OutPath: "deployments/k8s/base/kustomization.yaml",
		Include: func(c *ProjectConfig) bool { return c.IncludeK8s }},
	{TmplPath: "deployments/k8s/base/deployment.yaml.tmpl", OutPath: "deployments/k8s/base/deployment.yaml",
		Include: func(c *ProjectConfig) bool { return c.IncludeK8s }},
	{TmplPath: "deployments/k8s/base/service.yaml.tmpl", OutPath: "deployments/k8s/base/service.yaml",
		Include: func(c *ProjectConfig) bool { return c.IncludeK8s }},
	{TmplPath: "deployments/k8s/base/ingress.yaml.tmpl", OutPath: "deployments/k8s/base/ingress.yaml",
		Include: func(c *ProjectConfig) bool { return c.IncludeK8s }},
	{TmplPath: "deployments/k8s/base/configmap.yaml.tmpl", OutPath: "deployments/k8s/base/configmap.yaml",
		Include: func(c *ProjectConfig) bool { return c.IncludeK8s }},
	{TmplPath: "deployments/k8s/base/hpa.yaml.tmpl", OutPath: "deployments/k8s/base/hpa.yaml",
		Include: func(c *ProjectConfig) bool { return c.IncludeK8s }},
	{TmplPath: "deployments/k8s/base/networkpolicy.yaml.tmpl", OutPath: "deployments/k8s/base/networkpolicy.yaml",
		Include: func(c *ProjectConfig) bool { return c.IncludeK8s }},

	// --- SonarQube ---
	{TmplPath: "sonar-project.properties.tmpl", OutPath: "sonar-project.properties",
		Include: func(c *ProjectConfig) bool { return c.IncludeSonar }},

	// --- CI/CD ---
	{TmplPath: "dot-github/workflows/ci.yml.tmpl", OutPath: ".github/workflows/ci.yml",
		Include: func(c *ProjectConfig) bool { return c.IncludeGitHubCI }},
	{TmplPath: "dot-gitlab-ci.yml.tmpl", OutPath: ".gitlab-ci.yml",
		Include: func(c *ProjectConfig) bool { return c.IncludeGitLabCI }},

	// --- Shared infra (always included when any external service is selected) ---
	{TmplPath: "pkg/httpclient/client.go.tmpl", OutPath: "pkg/httpclient/client.go",
		Include: func(c *ProjectConfig) bool {
			return c.IncludeSendGrid || c.IncludeStripe || c.IncludeIceWarp || c.IncludeElasticsearch
		}},
	{TmplPath: "pkg/httpclient/service.go.tmpl", OutPath: "pkg/httpclient/service.go",
		Include: func(c *ProjectConfig) bool {
			return c.IncludeSendGrid || c.IncludeStripe || c.IncludeIceWarp || c.IncludeElasticsearch
		}},
	{TmplPath: "pkg/httpclient/codec.go.tmpl", OutPath: "pkg/httpclient/codec.go",
		Include: func(c *ProjectConfig) bool {
			return c.IncludeSendGrid || c.IncludeStripe || c.IncludeIceWarp || c.IncludeElasticsearch
		}},
	{TmplPath: "pkg/httpclient/token.go.tmpl", OutPath: "pkg/httpclient/token.go",
		Include: func(c *ProjectConfig) bool {
			return c.IncludeSendGrid || c.IncludeStripe || c.IncludeIceWarp || c.IncludeElasticsearch
		}},

	// --- Retry (needed by httpclient and broker) ---
	{TmplPath: "pkg/retry/retry.go.tmpl", OutPath: "pkg/retry/retry.go",
		Include: func(c *ProjectConfig) bool {
			return c.IncludeKafka || c.IncludeSendGrid || c.IncludeStripe || c.IncludeIceWarp || c.IncludeElasticsearch
		}},

	// --- Circuit breaker (needed by httpclient and broker) ---
	{TmplPath: "pkg/circuitbreaker/circuitbreaker.go.tmpl", OutPath: "pkg/circuitbreaker/circuitbreaker.go",
		Include: func(c *ProjectConfig) bool {
			return c.IncludeKafka || c.IncludeSendGrid || c.IncludeStripe || c.IncludeIceWarp || c.IncludeElasticsearch
		}},

	// --- Async (worker pool — always included, used by many features) ---
	{TmplPath: "pkg/async/worker.go.tmpl", OutPath: "pkg/async/worker.go"},
	{TmplPath: "pkg/async/parallel.go.tmpl", OutPath: "pkg/async/parallel.go"},
	{TmplPath: "pkg/async/context.go.tmpl", OutPath: "pkg/async/context.go"},

	// --- App services registry ---
	{TmplPath: "pkg/app/services.go.tmpl", OutPath: "pkg/app/services.go",
		Include: func(c *ProjectConfig) bool {
			return c.IncludeRedis || c.IncludeKafka || c.IncludeSendGrid || c.IncludeStripe ||
				c.IncludeIceWarp || c.IncludeFirebase || c.IncludeElasticsearch || c.IncludeEncryption
		}},
}
