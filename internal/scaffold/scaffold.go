package scaffold

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

//go:embed all:templates
var templateFS embed.FS

// Run scaffolds a new project in a directory named cfg.ProjectName under cwd.
func Run(cfg *ProjectConfig) error {
	return RunTo(cfg.ProjectName, cfg)
}

// RunTo scaffolds a new project into the specified output directory.
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

	// Process each file in manifest
	count := 0
	for _, rule := range manifest {
		if rule.Include != nil && !rule.Include(cfg) {
			continue
		}

		if err := processTemplate(cfg, rule, outputDir); err != nil {
			return fmt.Errorf("process %s: %w", rule.TmplPath, err)
		}
		count++
	}

	fmt.Printf("  Created %d files\n", count)

	// Run go mod tidy
	if err := runCommand(outputDir, "go", "mod", "tidy"); err != nil {
		fmt.Fprintf(os.Stderr, "  warning: go mod tidy failed: %v\n", err)
	} else {
		fmt.Println("  go mod tidy")
	}

	// Install binary
	binPath := "./cmd/" + cfg.ProjectName
	if err := runCommand(outputDir, "go", "install", binPath); err != nil {
		fmt.Fprintf(os.Stderr, "  warning: go install failed: %v\n", err)
	} else {
		fmt.Printf("  go install %s\n", binPath)
	}

	// Init git repo
	if err := runCommand(outputDir, "git", "init"); err != nil {
		fmt.Fprintf(os.Stderr, "  warning: git init failed: %v\n", err)
	} else {
		fmt.Println("  git init")
	}

	return nil
}

func processTemplate(cfg *ProjectConfig, rule FileRule, outputDir string) error {
	// Read template content from embedded FS
	tmplPath := "templates/" + rule.TmplPath
	content, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", tmplPath, err)
	}

	// Parse template
	tmpl, err := template.New(rule.TmplPath).
		Funcs(funcMap).
		Parse(string(content))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	// Resolve output path (dynamic or static)
	resolvedPath := rule.OutPath
	if rule.OutPathFunc != nil {
		resolvedPath = rule.OutPathFunc(cfg)
	}
	outPath := filepath.Join(outputDir, resolvedPath)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute template
	return tmpl.Execute(f, cfg)
}

func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
