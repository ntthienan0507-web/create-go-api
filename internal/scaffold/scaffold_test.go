package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func TestScaffold(t *testing.T) {
	tests := []struct {
		name         string
		config       ProjectConfig
		mustExist    []string
		mustNotExist []string
	}{
		{
			name: "sqlc_jwt_with_module",
			config: ProjectConfig{
				ModulePath:          "github.com/test/my-api",
				DBDriver:            "sqlc",
				AuthProvider:        "jwt",
				IncludeSampleModule: true,
			},
			mustExist: []string{
				"pkg/app/app.go",
				"pkg/config/config.go",
				"pkg/auth/provider.go",
				"pkg/auth/factory.go",
				"pkg/auth/jwt.go",
				"pkg/database/postgres.go",
				"pkg/database/store.go",
				"pkg/database/migration.go",
				"modules/user/controller.go",
				"modules/user/service.go",
				"modules/user/repository.go",
				"modules/user/repository_sqlc.go",
				"db/sqlc/db.go",
				"db/queries/users.sql",
				"sqlc.yaml",
				"go.mod",
				".env.example",
				".gitignore",
			},
			mustNotExist: []string{
				"pkg/auth/keycloak.go",
				"pkg/database/gorm.go",
				"modules/user/repository_gorm.go",
				"modules/user/repository_mongo.go",
				"pkg/broker",
				"pkg/ws",
				"pkg/cron",
				"pkg/crypto",
				"pkg/tracing",
				"deployments",
			},
		},
		{
			name: "gorm_keycloak_with_module",
			config: ProjectConfig{
				ModulePath:          "github.com/test/my-api",
				DBDriver:            "gorm",
				AuthProvider:        "keycloak",
				IncludeSampleModule: true,
			},
			mustExist: []string{
				"pkg/auth/keycloak.go",
				"pkg/database/gorm.go",
				"modules/user/repository_gorm.go",
				"modules/user/controller.go",
			},
			mustNotExist: []string{
				"pkg/auth/jwt.go",
				"pkg/database/postgres.go",
				"pkg/database/store.go",
				"modules/user/repository_sqlc.go",
				"db/sqlc/db.go",
				"sqlc.yaml",
			},
		},
		{
			name: "both_both_with_module_swagger",
			config: ProjectConfig{
				ModulePath:          "github.com/test/my-api",
				DBDriver:            "both",
				AuthProvider:        "both",
				IncludeSampleModule: true,
				IncludeSwagger:      true,
			},
			mustExist: []string{
				"pkg/auth/jwt.go",
				"pkg/auth/keycloak.go",
				"pkg/database/postgres.go",
				"pkg/database/store.go",
				"pkg/database/gorm.go",
				"modules/user/repository_sqlc.go",
				"modules/user/repository_gorm.go",
				"db/sqlc/db.go",
				"sqlc.yaml",
			},
			mustNotExist: []string{
				"modules/user/repository_mongo.go",
			},
		},
		{
			name: "sqlc_jwt_no_module",
			config: ProjectConfig{
				ModulePath:   "github.com/test/my-api",
				DBDriver:     "sqlc",
				AuthProvider: "jwt",
			},
			mustExist: []string{
				"pkg/app/app.go",
				"pkg/auth/jwt.go",
				"pkg/database/postgres.go",
				"go.mod",
			},
			mustNotExist: []string{
				"modules/user",
			},
		},
		{
			name: "gorm_with_kafka_redis",
			config: ProjectConfig{
				ModulePath:   "github.com/test/my-api",
				DBDriver:     "gorm",
				AuthProvider: "jwt",
				IncludeRedis: true,
				IncludeKafka: true,
			},
			mustExist: []string{
				"pkg/database/gorm.go",
				"pkg/database/redis.go",
				"pkg/broker/broker.go",
				"pkg/broker/producer.go",
				"pkg/broker/consumer.go",
				"pkg/broker/outbox.go",
				"pkg/cache",
				"pkg/retry/retry.go",
				"pkg/circuitbreaker/circuitbreaker.go",
			},
			mustNotExist: []string{
				"pkg/ws",
				"pkg/cron",
				"pkg/tracing",
				"deployments",
			},
		},
		{
			name: "full_stack",
			config: ProjectConfig{
				ModulePath:           "github.com/test/my-api",
				DBDriver:             "gorm",
				AuthProvider:         "jwt",
				IncludeSampleModule:  true,
				IncludeRedis:         true,
				IncludeKafka:         true,
				IncludeEncryption:    true,
				IncludeCron:          true,
				IncludeWebSocket:     true,
				IncludeOTEL:          true,
				IncludeSendGrid:      true,
				IncludeStripe:        true,
				IncludeFirebase:      true,
				IncludeElasticsearch: true,
				IncludeK8s:          true,
				IncludeSonar:        true,
				CIProvider:          "both",
			},
			mustExist: []string{
				"pkg/broker/broker.go",
				"pkg/crypto/crypto.go",
				"pkg/cron/cron.go",
				"pkg/ws/hub.go",
				"pkg/ws/client.go",
				"pkg/tracing/tracing.go",
				"pkg/external/sendgrid/client.go",
				"pkg/external/stripe/client.go",
				"pkg/external/firebase/client.go",
				"pkg/external/elasticsearch/client.go",
				"pkg/httpclient/client.go",
				"pkg/retry/retry.go",
				"pkg/circuitbreaker/circuitbreaker.go",
				"deployments/k8s/base/deployment.yaml",
				".github/workflows/ci.yml",
				".gitlab-ci.yml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outDir := filepath.Join(tmpDir, "my-api")

			tt.config.Derive()
			if err := RunTo(outDir, &tt.config); err != nil {
				t.Fatalf("RunTo failed: %v", err)
			}

			for _, f := range tt.mustExist {
				path := filepath.Join(outDir, f)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("expected file %s to exist", f)
				}
			}

			for _, f := range tt.mustNotExist {
				path := filepath.Join(outDir, f)
				if _, err := os.Stat(path); err == nil {
					t.Errorf("expected file %s to NOT exist", f)
				}
			}

			// Verify no source module path leaked into output
			filepath.Walk(outDir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".mod") {
					return nil
				}
				data, _ := os.ReadFile(path)
				if strings.Contains(string(data), sourceModulePath) {
					rel, _ := filepath.Rel(outDir, path)
					t.Errorf("file %s still contains source module path %q", rel, sourceModulePath)
				}
				return nil
			})
		})
	}
}

// TestAllTemplatesParse is skipped — templates are now cloned from gostack-kit at runtime,
// not embedded. Template validity is verified by the gostack-kit CI pipeline.

func parseTemplate(name, content string) (interface{}, error) {
	_, err := template.New(name).Funcs(funcMap).Parse(content)
	return nil, err
}
