package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// sourceRepo is the gostack-kit template repository.
	sourceRepo = "https://github.com/ntthienan0507-web/gostack-kit.git"
	// sourceModulePath is the module path in gostack-kit that gets replaced.
	sourceModulePath = "github.com/ntthienan0507-web/gostack-kit"
)

// Run scaffolds a new project in a directory named cfg.ProjectName under cwd.
func Run(cfg *ProjectConfig) error {
	return RunTo(cfg.ProjectName, cfg)
}

// RunTo scaffolds a new project into the specified output directory.
// It clones gostack-kit, filters files based on config, and replaces the module path.
func RunTo(outputDir string, cfg *ProjectConfig) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}

	// Check directory is empty
	entries, _ := os.ReadDir(outputDir)
	if len(entries) > 0 {
		return fmt.Errorf("directory %q is not empty", outputDir)
	}

	// Step 1: Clone gostack-kit into temp dir
	fmt.Println("  Cloning gostack-kit...")
	tmpDir, err := os.MkdirTemp("", "gostack-kit-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := runCommand(".", "git", "clone", "--depth=1", sourceRepo, tmpDir); err != nil {
		return fmt.Errorf("clone gostack-kit: %w", err)
	}

	// Remove .git from clone
	os.RemoveAll(filepath.Join(tmpDir, ".git"))

	// Step 2: Filter based on manifest rules (just used for Include logic, not template processing)
	fmt.Println("  Selecting files based on config...")
	count := 0

	// Step 3: Copy remaining core files that aren't in manifest but always needed
	coreCopies := []struct{ src, dst string }{
		{"go.mod", "go.mod"},
		{"go.sum", "go.sum"},
		{"Makefile", "Makefile"},
		{"Dockerfile", "Dockerfile"},
		{".env.example", ".env.example"},
		{".gitignore", ".gitignore"},
		{".editorconfig", ".editorconfig"},
	}
	for _, c := range coreCopies {
		src := filepath.Join(tmpDir, c.src)
		if _, err := os.Stat(src); err == nil {
			dst := filepath.Join(outputDir, c.dst)
			copyFile(src, dst, cfg.ModulePath)
			count++
		}
	}

	// Step 4: Copy directories that always exist
	alwaysDirs := []string{
		"pkg/config", "pkg/logger", "pkg/response", "pkg/apperror",
		"pkg/auth", "pkg/middleware", "pkg/database", "pkg/app",
		"pkg/async", "db/migrations/00001_create_users.sql",
	}
	for _, d := range alwaysDirs {
		src := filepath.Join(tmpDir, d)
		if info, err := os.Stat(src); err == nil {
			dst := filepath.Join(outputDir, d)
			if info.IsDir() {
				n, _ := copyDir(src, dst, cfg.ModulePath)
				count += n
			} else {
				copyFile(src, dst, cfg.ModulePath)
				count++
			}
		}
	}

	// Step 4b: Remove DB-driver-exclusive files
	if !cfg.IncludeSQLC {
		for _, f := range []string{"pkg/database/postgres.go", "pkg/database/store.go"} {
			os.Remove(filepath.Join(outputDir, f))
		}
	}
	if !cfg.IncludeGORM {
		os.Remove(filepath.Join(outputDir, "pkg/database/gorm.go"))
		os.Remove(filepath.Join(outputDir, "pkg/database/transaction.go"))
		os.Remove(filepath.Join(outputDir, "pkg/database/transaction_test.go"))
	}

	// Step 4c: Remove auth-exclusive files
	if !cfg.IncludeJWT {
		os.Remove(filepath.Join(outputDir, "pkg/auth/jwt.go"))
		os.Remove(filepath.Join(outputDir, "pkg/auth/jwt_test.go"))
	}
	if !cfg.IncludeKeycloak {
		os.Remove(filepath.Join(outputDir, "pkg/auth/keycloak.go"))
	}

	// Step 5: Copy sample module if enabled
	if cfg.IncludeSampleModule {
		src := filepath.Join(tmpDir, "modules/user")
		dst := filepath.Join(outputDir, "modules/user")
		if _, err := os.Stat(src); err == nil {
			n, _ := copyDir(src, dst, cfg.ModulePath)
			count += n

			// Remove db-specific repos based on driver choice
			if !cfg.IncludeSQLC {
				os.Remove(filepath.Join(dst, "repository_sqlc.go"))
			}
			if !cfg.IncludeGORM {
				os.Remove(filepath.Join(dst, "repository_gorm.go"))
			}
			// Remove mongo repo if not selected (currently no mongo option)
			os.Remove(filepath.Join(dst, "repository_mongo.go"))

			// Remove events.go if no kafka
			if !cfg.IncludeKafka {
				os.Remove(filepath.Join(dst, "events.go"))
			}
		}
	}

	// Step 6: Copy SQLC files if needed
	if cfg.IncludeSQLC {
		for _, d := range []string{"db/sqlc", "db/queries"} {
			src := filepath.Join(tmpDir, d)
			dst := filepath.Join(outputDir, d)
			if _, err := os.Stat(src); err == nil {
				copyDir(src, dst, cfg.ModulePath)
			}
		}
		src := filepath.Join(tmpDir, "sqlc.yaml")
		if _, err := os.Stat(src); err == nil {
			copyFile(src, filepath.Join(outputDir, "sqlc.yaml"), cfg.ModulePath)
		}
	}

	// Step 7: Copy optional packages based on config
	optionalPackages := []struct {
		enabled bool
		paths   []string
	}{
		{cfg.IncludeRedis, []string{"pkg/cache"}},
		{cfg.IncludeKafka, []string{"pkg/broker"}},
		{cfg.IncludeEncryption, []string{"pkg/crypto"}},
		{cfg.IncludeCron, []string{"pkg/cron"}},
		{cfg.IncludeWebSocket, []string{"pkg/ws"}},
		{cfg.IncludeOTEL, []string{"pkg/tracing"}},
		{cfg.IncludeSendGrid, []string{"pkg/external/sendgrid"}},
		{cfg.IncludeStripe, []string{"pkg/external/stripe"}},
		{cfg.IncludeIceWarp, []string{"pkg/external/icewarp"}},
		{cfg.IncludeFirebase, []string{"pkg/external/firebase"}},
		{cfg.IncludeElasticsearch, []string{"pkg/external/elasticsearch"}},
		{cfg.IncludeK8s, []string{"deployments"}},
		{cfg.IncludeSonar, []string{"sonar-project.properties", "scripts"}},
		{cfg.IncludeGitHubCI, []string{".github"}},
		{cfg.IncludeGitLabCI, []string{".gitlab-ci.yml"}},
		// Shared infra needed by external services
		{cfg.IncludeSendGrid || cfg.IncludeStripe || cfg.IncludeIceWarp || cfg.IncludeElasticsearch,
			[]string{"pkg/httpclient"}},
		{cfg.IncludeKafka || cfg.IncludeSendGrid || cfg.IncludeStripe || cfg.IncludeIceWarp || cfg.IncludeElasticsearch,
			[]string{"pkg/retry", "pkg/circuitbreaker"}},
	}

	for _, opt := range optionalPackages {
		if !opt.enabled {
			continue
		}
		for _, p := range opt.paths {
			src := filepath.Join(tmpDir, p)
			info, err := os.Stat(src)
			if err != nil {
				continue
			}
			dst := filepath.Join(outputDir, p)
			if info.IsDir() {
				n, _ := copyDir(src, dst, cfg.ModulePath)
				count += n
			} else {
				copyFile(src, dst, cfg.ModulePath)
				count++
			}
		}
	}

	fmt.Printf("  Created %d files\n", count)

	// Step 8: Rewrite go.mod module path
	rewriteGoMod(outputDir, cfg)

	// Step 9: go mod tidy
	if err := runCommand(outputDir, "go", "mod", "tidy"); err != nil {
		fmt.Fprintf(os.Stderr, "  warning: go mod tidy failed: %v\n", err)
	} else {
		fmt.Println("  go mod tidy")
	}

	// Step 10: git init
	if err := runCommand(outputDir, "git", "init"); err != nil {
		fmt.Fprintf(os.Stderr, "  warning: git init failed: %v\n", err)
	} else {
		fmt.Println("  git init")
	}

	return nil
}

// copyFile copies a single file, replacing the source module path with the target module path.
func copyFile(src, dst, targetModulePath string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Replace module path in Go files, YAML, properties, etc.
	content := string(data)
	if shouldReplace(src) {
		content = strings.ReplaceAll(content, sourceModulePath, targetModulePath)
	}

	return os.WriteFile(dst, []byte(content), 0o644)
}

// copyDir recursively copies a directory, replacing module paths.
func copyDir(src, dst, targetModulePath string) (int, error) {
	count := 0
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip test files, binary files, .git
		base := filepath.Base(path)
		if base == ".git" || base == ".DS_Store" {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		if err := copyFile(path, target, targetModulePath); err != nil {
			return err
		}
		count++
		return nil
	})
	return count, err
}

// shouldReplace returns true if the file content should have module path replaced.
func shouldReplace(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".go", ".mod", ".yaml", ".yml", ".properties", ".json", ".toml", ".md":
		return true
	}
	base := filepath.Base(path)
	switch base {
	case "Makefile", "Dockerfile", ".env.example":
		return true
	}
	return false
}

// rewriteGoMod rewrites the go.mod module path.
func rewriteGoMod(outputDir string, cfg *ProjectConfig) {
	modPath := filepath.Join(outputDir, "go.mod")
	data, err := os.ReadFile(modPath)
	if err != nil {
		return
	}

	content := strings.Replace(string(data), sourceModulePath, cfg.ModulePath, 1)
	os.WriteFile(modPath, []byte(content), 0o644)
	fmt.Println("  go.mod → " + cfg.ModulePath)
}

func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
