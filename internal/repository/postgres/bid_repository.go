package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
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

func (r *bidRepository) GetBidById(ctx context.Context, id string) (*model.Bid, error) {
	stmt, err := r.db.PrepareContext(ctx, `
	SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
	FROM bid
	WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var bid model.Bid
	err = stmt.QueryRowContext(ctx, id).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.CreatorUsername,
		&bid.Version,
		&bid.CreatedAt,
		&bid.UpdatedAt,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting bid", slog.Any("error", err))
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
		return nil, err
	}
	return &bid, nil
}

func (r *bidRepository) GetBidByUsername(ctx context.Context, limit int, offset int, username string) ([]model.Bid, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE creator_username = $1
		LIMIT $2 OFFSET $3
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username, limit, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
		return nil, err
	}
	defer rows.Close()

	var bids []model.Bid
	for rows.Next() {
		bid := model.Bid{}
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.CreatorUsername, &bid.Version, &bid.CreatedAt, &bid.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (r *bidRepository) GetTenderBids(ctx context.Context, tenderID string, limit int, offset int, username string) ([]model.Bid, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE tender_id = $1
		LIMIT $2 OFFSET $3
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, tenderID, limit, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting bids", slog.Any("error", err))
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
		return nil, err
	}
	defer rows.Close()

	var bids []model.Bid
	for rows.Next() {
		bid := model.Bid{}
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID, &bid.CreatorUsername, &bid.Version, &bid.CreatedAt, &bid.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (r *bidRepository) GetBidStatus(ctx context.Context, bidID string) (model.BidStatus, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT status
		FROM bid
		WHERE id = $1
	`)
	if err != nil {
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var status model.BidStatus
	err = stmt.QueryRowContext(ctx, bidID).Scan(&status)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting bid status", slog.Any("error", err))
			return "", fmt.Errorf("failed to execute query: %w", err)
		}
		return "", err
	}

	return status, nil
}

func (r *bidRepository) UpdateBid(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE bid
		SET name = $1, description = $2, status = $3, tender_id = $4, author_type = $5, author_id = $6, creator_username = $7, version = $8, created_at = $9, updated_at = $10
		WHERE id = $11
		RETURNING id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	bid.Version++
	bid.UpdatedAt = time.Now()

	err = stmt.QueryRowContext(ctx,
		bid.Name,
		bid.Description,
		bid.Status,
		bid.TenderID,
		bid.AuthorType,
		bid.AuthorID,
		bid.CreatorUsername,
		bid.Version,
		bid.CreatedAt,
		bid.UpdatedAt,
		bid.ID,
	).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.CreatorUsername,
		&bid.Version,
		&bid.CreatedAt,
		&bid.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return bid, nil
}
