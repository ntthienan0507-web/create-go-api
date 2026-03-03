package scaffold

import (
	"io/fs"
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
				"cmd/server/main.go",
				"cmd/cli/main.go",
				"cmd/dbtest/main.go",
				"internal/app/app.go",
				"internal/config/config.go",
				"internal/auth/provider.go",
				"internal/auth/factory.go",
				"internal/auth/jwt.go",
				"internal/database/postgres.go",
				"internal/database/store.go",
				"internal/database/migration.go",
				"internal/module/user/handler.go",
				"internal/module/user/service.go",
				"internal/module/user/repository.go",
				"internal/module/user/repository_sqlc.go",
				"db/sqlc/db.go",
				"db/queries/users.sql",
				"sqlc.yaml",
				"go.mod",
				".env.example",
				".gitignore",
			},
			mustNotExist: []string{
				"internal/auth/keycloak.go",
				"internal/database/gorm.go",
				"internal/module/user/repository_gorm.go",
				"docs/.gitkeep",
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
				"internal/auth/keycloak.go",
				"internal/database/gorm.go",
				"internal/module/user/repository_gorm.go",
				"internal/module/user/handler.go",
			},
			mustNotExist: []string{
				"internal/auth/jwt.go",
				"internal/database/postgres.go",
				"internal/database/store.go",
				"internal/module/user/repository_sqlc.go",
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
				"internal/auth/jwt.go",
				"internal/auth/keycloak.go",
				"internal/database/postgres.go",
				"internal/database/store.go",
				"internal/database/gorm.go",
				"internal/module/user/repository_sqlc.go",
				"internal/module/user/repository_gorm.go",
				"db/sqlc/db.go",
				"sqlc.yaml",
				"docs/.gitkeep",
			},
			mustNotExist: []string{},
		},
		{
			name: "sqlc_jwt_no_module",
			config: ProjectConfig{
				ModulePath:          "github.com/test/my-api",
				DBDriver:            "sqlc",
				AuthProvider:        "jwt",
				IncludeSampleModule: false,
			},
			mustExist: []string{
				"internal/auth/jwt.go",
				"internal/database/postgres.go",
				"internal/database/store.go",
				"db/sqlc/db.go",
			},
			mustNotExist: []string{
				"internal/module/user/handler.go",
				"internal/module/user/service.go",
				"internal/module/user/repository.go",
				"internal/module/user/repository_sqlc.go",
				"internal/module/user/repository_gorm.go",
				"db/queries/users.sql",
				"db/sqlc/users.sql.go",
				"db/sqlc/models.go",
				"db/sqlc/querier.go",
			},
		},
		{
			name: "gorm_jwt_no_module",
			config: ProjectConfig{
				ModulePath:          "github.com/test/my-api",
				DBDriver:            "gorm",
				AuthProvider:        "jwt",
				IncludeSampleModule: false,
			},
			mustExist: []string{
				"internal/auth/jwt.go",
				"internal/database/gorm.go",
			},
			mustNotExist: []string{
				"internal/database/postgres.go",
				"internal/database/store.go",
				"internal/module/user/handler.go",
				"db/sqlc/db.go",
			},
		},
		{
			name: "both_keycloak_no_module",
			config: ProjectConfig{
				ModulePath:          "github.com/test/my-api",
				DBDriver:            "both",
				AuthProvider:        "keycloak",
				IncludeSampleModule: false,
			},
			mustExist: []string{
				"internal/auth/keycloak.go",
				"internal/database/postgres.go",
				"internal/database/gorm.go",
			},
			mustNotExist: []string{
				"internal/auth/jwt.go",
				"internal/module/user/handler.go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			// Use a subdirectory so it's empty
			outDir := filepath.Join(dir, "project")

			tt.config.Derive()
			err := RunTo(outDir, &tt.config)
			if err != nil {
				t.Fatalf("RunTo failed: %v", err)
			}

			// Check expected files exist
			for _, f := range tt.mustExist {
				p := filepath.Join(outDir, f)
				if _, err := os.Stat(p); os.IsNotExist(err) {
					t.Errorf("expected file %s to exist", f)
				}
			}

			// Check excluded files don't exist
			for _, f := range tt.mustNotExist {
				p := filepath.Join(outDir, f)
				if _, err := os.Stat(p); err == nil {
					t.Errorf("expected file %s to NOT exist", f)
				}
			}

			// Walk all files and check for unresolved template syntax
			filepath.WalkDir(outDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				content, readErr := os.ReadFile(path)
				if readErr != nil {
					return nil
				}
				s := string(content)
				rel, _ := filepath.Rel(outDir, path)

				if strings.Contains(s, "github.com/ntthienan0507-web/go-api-template") {
					t.Errorf("file %s still contains old module path", rel)
				}
				if strings.Contains(s, "{{") && strings.Contains(s, "}}") {
					t.Errorf("file %s contains unresolved template syntax", rel)
				}
				return nil
			})

			// Check go.mod has correct module path
			goModPath := filepath.Join(outDir, "go.mod")
			goModContent, err := os.ReadFile(goModPath)
			if err != nil {
				t.Fatalf("read go.mod: %v", err)
			}
			if !strings.Contains(string(goModContent), tt.config.ModulePath) {
				t.Errorf("go.mod does not contain module path %q", tt.config.ModulePath)
			}
		})
	}
}

func TestAllTemplatesParse(t *testing.T) {
	// Validate every .tmpl file can be parsed without errors
	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		content, readErr := templateFS.ReadFile(path)
		if readErr != nil {
			t.Errorf("read %s: %v", path, readErr)
			return nil
		}

		_, parseErr := parseTemplate(path, string(content))
		if parseErr != nil {
			t.Errorf("parse %s: %v", path, parseErr)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk templates: %v", err)
	}
}

func parseTemplate(name, content string) (interface{}, error) {
	_, err := template.New(name).Funcs(funcMap).Parse(content)
	return nil, err
}
