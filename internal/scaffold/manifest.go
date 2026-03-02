package scaffold

// FileRule defines when a template file should be included.
type FileRule struct {
	TmplPath string                    // path inside templates/ FS
	OutPath  string                    // output path (minus .tmpl)
	Include  func(*ProjectConfig) bool // nil = always include
}

// manifest lists every template file and its inclusion predicate.
var manifest = []FileRule{
	// --- Always included (core skeleton) ---
	{TmplPath: "cmd/server/main.go.tmpl", OutPath: "cmd/server/main.go"},
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
}
