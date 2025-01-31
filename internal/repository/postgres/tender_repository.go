package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"
)

type tenderRepository struct {
	db *sql.DB
	logger *slog.Logger
}


func NewTenderRepository(db *sql.DB, logger *slog.Logger) repository.TenderRepository {
	return &tenderRepository{
		db: db,
		logger: logger,
	}
}


func (r *tenderRepository) CreateTender(ctx context.Context, tender *model.Tender) (*model.Tender, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
    if err := tx.Rollback(); err != nil {
        if err != sql.ErrTxDone && err != sql.ErrConnDone { 
            r.logger.ErrorContext(ctx, "Error rolling back transaction", slog.Any("error", err))
        }
    }
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO tender (id, name, description, service_type, organization_id, creator_username, status, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for creating tender: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,
		tender.ID,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.OrganizationID,
		tender.CreatorUsername,
		tender.Status,
		tender.Version,
	)

	if err := row.Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.OrganizationID,
		&tender.CreatorUsername,
		&tender.Status,
		&tender.Version,
		&tender.CreatedAt,
		&tender.UpdatedAt,
	); err != nil {
		r.logger.ErrorContext(ctx, "Error creating tender and scanning result", slog.Any("error", err))
		return nil, fmt.Errorf("failed to execute query and scan result: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tender, nil
}

func (r *tenderRepository) GetTenders(ctx context.Context, limit int, offset int, serviceTypes []model.TenderServiceType) ([]model.Tender, error) {

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		FROM tender
		WHERE service_type = ANY($1)
		LIMIT $2 OFFSET $3
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for getting tenders: %w", err)
	}
	defer stmt.Close()

	serviceTypeStrings := make([]string, len(serviceTypes))
	for i, st := range serviceTypes {
		serviceTypeStrings[i] = string(st)
	}

	rows, err := stmt.QueryContext(ctx, pq.Array(serviceTypeStrings), limit, offset) 
	if err != nil {
		r.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return nil, fmt.Errorf("failed to execute query for getting tenders: %w", err)
	}
	
	defer rows.Close()

	var tenders []model.Tender
	for rows.Next() {
		tender := model.Tender{}
		if err := rows.Scan(
			&tender.ID,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.OrganizationID,
			&tender.CreatorUsername,
			&tender.Status,
			&tender.Version,
			&tender.CreatedAt,
			&tender.UpdatedAt,
		); err != nil {
			r.logger.ErrorContext(ctx, "Error scanning tender", slog.Any("error", err))
			return nil, fmt.Errorf("failed to scan tender: %w", err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}


func (r *tenderRepository) GetTenderById(ctx context.Context, id string) (*model.Tender, error) {

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		FROM tender
		WHERE id = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for getting tender by id: %w", err)
	}
	defer stmt.Close()

	tender := model.Tender{}
	err = stmt.QueryRowContext(ctx, id).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.OrganizationID,
		&tender.CreatorUsername,
		&tender.Status,	
		&tender.Version,
		&tender.CreatedAt,
		&tender.UpdatedAt,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting tender by id", slog.Any("error", err))
			return nil, fmt.Errorf("failed to execute query for getting tender by id: %w", err)
		}
		return nil, err	
	}

	return &tender, nil
}



func (r *tenderRepository) GetTenderByUsername(ctx context.Context, limit int, offset int, username string) ([]model.Tender, error) {

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		FROM tender
		WHERE creator_username = $1
		LIMIT $2 OFFSET $3
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for getting tenders: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username, limit, offset) 
	if err != nil {
		r.logger.ErrorContext(ctx, "Error getting tenders", slog.Any("error", err))
		return nil, fmt.Errorf("failed to execute query for getting tenders: %w", err)
	}
	
	defer rows.Close()

	var tenders []model.Tender
	for rows.Next() {
		tender := model.Tender{}
		if err := rows.Scan(
			&tender.ID,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.OrganizationID,			
			&tender.CreatorUsername,
			&tender.Status,
			&tender.Version,
			&tender.CreatedAt,
			&tender.UpdatedAt,
		); err != nil {
			r.logger.ErrorContext(ctx, "Error scanning tender", slog.Any("error", err))
			return nil, fmt.Errorf("failed to scan tender: %w", err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}


func (r *tenderRepository) IsUserResponsibleForTender(ctx context.Context, tenderID string, username string) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM organization_responsible orr
			JOIN tender t ON t.organization_id = orr.organization_id
			JOIN employee e ON e.id = orr.user_id
			WHERE t.id = $1 AND e.username = $2
		)
	`)
	if err != nil {
		return false, fmt.Errorf("error preparing statement for checking user is responsible: %w", err)
	}

	defer stmt.Close()

	

	var exists bool

	err = stmt.QueryRowContext(ctx, tenderID, username).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("error checking user is responsible: %w", err)
	}

	return exists, nil
}



func (r *tenderRepository) UpdateTender(ctx context.Context, tender *model.Tender) (*model.Tender, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
    if err := tx.Rollback(); err != nil {
       
        if err != sql.ErrTxDone && err != sql.ErrConnDone { 
            r.logger.ErrorContext(ctx, "Error rolling back transaction", slog.Any("error", err))
        }
    }
	}()
	
	stmt, err := tx.PrepareContext(ctx, `
		UPDATE tender
		SET name = $2, description = $3, service_type = $4, organization_id = $5, creator_username = $6, status = $7, version = $8, updated_at = $9
		WHERE id = $1
		RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for updating tender: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,
		tender.ID,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.OrganizationID,
		tender.CreatorUsername,
		tender.Status,
		tender.Version+1,
		time.Now(),
	)

	var updatedTender model.Tender
	if err := row.Scan(
		&updatedTender.ID,
		&updatedTender.Name,
		&updatedTender.Description,
		&updatedTender.ServiceType,
		&updatedTender.OrganizationID,
		&updatedTender.CreatorUsername,
		&updatedTender.Status,
		&updatedTender.Version,
		&updatedTender.CreatedAt,
		&updatedTender.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to scan updated tender: %w", err)
	}	

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &updatedTender, nil
	}