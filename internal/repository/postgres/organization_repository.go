package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"database/sql"
	"log/slog"
)

type organizationRepository struct {
	db *sql.DB
	logger *slog.Logger
}

func NewOrganizationRepository(db *sql.DB,logger *slog.Logger) repository.OrganizationRepository {
	return &organizationRepository{
		db: db,
		logger: logger,
	}
}