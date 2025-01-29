package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type bidRepository struct {
	db     *sql.DB
	logger *slog.Logger
}


func NewBidRepository(db *sql.DB, logger *slog.Logger) repository.BidRepository {
	return &bidRepository{
		db:     db,
		logger: logger,
	}
}


func (r *bidRepository) CreateBid(ctx context.Context, bidRequest *model.Bid) (*model.Bid, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() 

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO bid (id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var bid model.Bid
	err = stmt.QueryRowContext(ctx, bidRequest.ID, bidRequest.Name, bidRequest.Description, bidRequest.Status, bidRequest.TenderID, bidRequest.AuthorType, bidRequest.AuthorID, bidRequest.CreatorUsername, bidRequest.Version, bidRequest.CreatedAt, bidRequest.UpdatedAt).Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.CreatorUsername, &bid.Version, &bid.CreatedAt, &bid.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &bid, nil
}