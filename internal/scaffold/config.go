package scaffold

import (
	"path"
	"strings"
)

// ProjectConfig holds all user choices passed as template data.
type ProjectConfig struct {
	// Core
	ModulePath  string // e.g. "github.com/myorg/my-api"
	ProjectName string // e.g. "my-api" (last segment of ModulePath)

	// Database
	DBDriver    string // "sqlc" | "gorm" | "both"
	IncludeSQLC bool   // derived
	IncludeGORM bool   // derived

	// Auth
	AuthProvider    string // "jwt" | "keycloak" | "both"
	IncludeJWT      bool   // derived
	IncludeKeycloak bool   // derived

	// Features
	IncludeSampleModule bool
	IncludeSwagger      bool

	// Infrastructure
	IncludeRedis      bool
	IncludeKafka      bool
	IncludeEncryption bool

	// Extra features
	IncludeCron      bool
	IncludeWebSocket bool
	IncludeOTEL      bool

	// External services
	IncludeSendGrid      bool
	IncludeStripe        bool
	IncludeIceWarp       bool
	IncludeFirebase      bool
	IncludeElasticsearch bool

	// DevOps
	IncludeK8s      bool
	IncludeSonar    bool
	CIProvider      string // "github" | "gitlab" | "both" | ""
	IncludeGitHubCI bool   // derived
	IncludeGitLabCI bool   // derived

	// Server
	ServerPort int    // default 8080
	DBName     string // default = ProjectName with hyphens → underscores
}

// Derive sets computed boolean fields from string choices and derives
// ProjectName / DBName from ModulePath if not already set.
func (c *ProjectConfig) Derive() {
	c.IncludeSQLC = c.DBDriver == "sqlc" || c.DBDriver == "both"
	c.IncludeGORM = c.DBDriver == "gorm" || c.DBDriver == "both"
	c.IncludeJWT = c.AuthProvider == "jwt" || c.AuthProvider == "both"
	c.IncludeKeycloak = c.AuthProvider == "keycloak" || c.AuthProvider == "both"

	c.IncludeGitHubCI = c.CIProvider == "github" || c.CIProvider == "both"
	c.IncludeGitLabCI = c.CIProvider == "gitlab" || c.CIProvider == "both"

	if c.ProjectName == "" {
		c.ProjectName = path.Base(c.ModulePath)
	}
	if c.DBName == "" {
		c.DBName = strings.ReplaceAll(c.ProjectName, "-", "_")
	}
	if c.ServerPort == 0 {
		c.ServerPort = 8080
	}
}
