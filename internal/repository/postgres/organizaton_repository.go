package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type organizationRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewOrganizationRepository(db *sql.DB, logger *slog.Logger) repository.OrganizationRepository {
	return &organizationRepository{db: db, logger: logger}
}

func (r *organizationRepository) GetOrganizationById(ctx context.Context, id string) (*model.Organization, error) {

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, name, description,type,created_at, updated_at
		FROM organization
		WHERE id = $1	
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for getting organization by id: %w", err)
	}
	defer stmt.Close()

	organization := model.Organization{}
	err = stmt.QueryRowContext(ctx, id).Scan(
		&organization.Id,
		&organization.Name,
		&organization.Description,
		&organization.Type,
		&organization.Created_at,	
		&organization.Updated_at,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting organization by id", slog.Any("error", err))
			return nil, fmt.Errorf("failed to execute query for getting organization by id: %w", err)
		}
		return nil, nil
	}

	return &organization, nil
}

func (r *organizationRepository) IsUserResponsibleForOrganization(ctx context.Context, organizationID string, username string) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM organization_responsible orr
			JOIN employee e ON e.id = orr.user_id
			WHERE orr.organization_id = $1 AND e.username = $2
		)
	`)
	if err != nil {
		return false, fmt.Errorf("error preparing statement for checking user is responsible: %w", err)
	}

	defer stmt.Close()

	var exists bool

	err = stmt.QueryRowContext(ctx, organizationID, username).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("error checking user is responsible: %w", err)
	}

	return exists, nil
}