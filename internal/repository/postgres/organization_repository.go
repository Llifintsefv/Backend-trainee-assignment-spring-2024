package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"database/sql"
)

type organizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) repository.OrganizationRepository {
	return &organizationRepository{
		db: db,
	}
}