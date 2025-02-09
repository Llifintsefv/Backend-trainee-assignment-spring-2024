package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
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
	defer func() {
		if err := tx.Rollback(); err != nil {
			if err != sql.ErrTxDone && err != sql.ErrConnDone {
				r.logger.ErrorContext(ctx, "Error rolling back transaction", slog.Any("error", err))
			}
		}
	}()

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

	oldVersion := bid.Version
	bid.Version++
	bid.UpdatedAt = time.Now()

	stmt, err := tx.PrepareContext(ctx, `
        UPDATE bid 
        SET name = $1, description = $2, status = $3, tender_id = $4, 
            author_type = $5, author_id = $6, creator_username = $7, 
            version = $8, created_at = $9, updated_at = $10 
        WHERE id = $11
        RETURNING id, name, description, status, tender_id, author_type, 
            author_id, creator_username, version, created_at, updated_at
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	updatedBid := new(model.Bid)
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
		&updatedBid.ID,
		&updatedBid.Name,
		&updatedBid.Description,
		&updatedBid.Status,
		&updatedBid.TenderID,
		&updatedBid.AuthorType,
		&updatedBid.AuthorID,
		&updatedBid.CreatorUsername,
		&updatedBid.Version,
		&updatedBid.CreatedAt,
		&updatedBid.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update bid: %w", err)
	}

	historyStmt, err := tx.PrepareContext(ctx, `
        INSERT INTO bid_history (
            id, bid_id, name, description, status, tender_id, 
            author_type, author_id, creator_username, 
            version, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare history statement: %w", err)
	}
	defer historyStmt.Close()

	_, err = historyStmt.ExecContext(ctx,
		uuid.New().String(),
		bid.ID,
		bid.Name,
		bid.Description,
		bid.Status,
		bid.TenderID,
		bid.AuthorType,
		bid.AuthorID,
		bid.CreatorUsername,
		oldVersion,
		bid.CreatedAt,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert bid history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedBid, nil
}

func (r *bidRepository) RollbackBidVersion(ctx context.Context, bidID string, version int) (*model.Bid, error) {
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

	var historyBid model.Bid
	err = tx.QueryRowContext(ctx, `
		SELECT id, name, description, status, tender_id, author_type, 
			author_id, creator_username, version, created_at, updated_at
		FROM bid_history
		WHERE bid_id = $1 AND version = $2
	`, bidID, version).Scan(
		&historyBid.ID,
		&historyBid.Name,
		&historyBid.Description,
		&historyBid.Status,
		&historyBid.TenderID,
		&historyBid.AuthorType,
		&historyBid.AuthorID,
		&historyBid.CreatorUsername,
		&historyBid.Version,
		&historyBid.CreatedAt,
		&historyBid.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query bid history: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE bid
		SET name = $1, description = $2, status = $3, tender_id = $4, 
			author_type = $5, author_id = $6, creator_username = $7, 
			version = $8, created_at = $9, updated_at = $10
		WHERE id = $11
		RETURNING id, name, description, status, tender_id, author_type, 
			author_id, creator_username, version, created_at, updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,
		historyBid.Name,
		historyBid.Description,
		historyBid.Status,
		historyBid.TenderID,
		historyBid.AuthorType,
		historyBid.AuthorID,
		historyBid.CreatorUsername,
		historyBid.Version,
		historyBid.CreatedAt,
		time.Now(),
		bidID,
	)

	var updatedBid model.Bid
	err = row.Scan(
		&updatedBid.ID,
		&updatedBid.Name,
		&updatedBid.Description,
		&updatedBid.Status,
		&updatedBid.TenderID,
		&updatedBid.AuthorType,
		&updatedBid.AuthorID,
		&updatedBid.CreatorUsername,
		&updatedBid.Version,
		&updatedBid.CreatedAt,
		&updatedBid.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update bid: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &updatedBid, nil
}

func (r *bidRepository) AddBidFeedback(ctx context.Context, bidID string, username string, review string) (*model.Bid, error) {
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
		INSERT INTO bid_feedback (id, bid_id, username, review, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uuid.New().String(), bidID, username, review, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert bid feedback: %w", err)
	}

	stmt, err = tx.PrepareContext(ctx, `
		UPDATE bid
		SET version = version + 1
		WHERE id = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bidID)
	if err != nil {
		return nil, fmt.Errorf("failed to update bid version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetBidById(ctx, bidID)

}
