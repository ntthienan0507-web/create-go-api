package main

import (
	"fmt"
	"os"

	"github.com/ntthienan0507-web/create-go-api/internal/prompt"
	"github.com/ntthienan0507-web/create-go-api/internal/scaffold"
)

const version = "0.1.0"

func main() {
	fmt.Println("create-go-api v" + version)
	fmt.Println()

	// Collect user input via interactive prompts
	cfg, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Scaffold project
	fmt.Printf("\nScaffolding %s...\n", cfg.ProjectName)
	if err := scaffold.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print next steps
	fmt.Println()
	fmt.Println("Done! Next steps:")
	fmt.Printf("  cd %s\n", cfg.ProjectName)
	fmt.Println("  cp .env.example .env   # edit with your DB credentials")
	fmt.Println("  make dev")
	fmt.Println()
	if cfg.IncludeSQLC {
		fmt.Println("SQLC tip: run 'make sqlc' after editing db/queries/*.sql")
	}
}
