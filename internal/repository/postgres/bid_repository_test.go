package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestBid(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository.BidRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	logger := slog.Default()
	repo := NewBidRepository(db, logger)

	return db, mock, repo
}

func TestCreateBid(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidRequest := &model.Bid{
			ID:              uuid.New().String(),
			Name:            "Test Bid",
			Description:     "Test Description",
			Status:          "open",
			TenderID:        uuid.New().String(),
			AuthorType:      "user",
			AuthorID:        uuid.New().String(),
			CreatorUsername: "testuser",
			Version:         1,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`
		INSERT INTO bid (id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
	`)).ExpectQuery().WithArgs(
			bidRequest.ID, bidRequest.Name, bidRequest.Description, bidRequest.Status, bidRequest.TenderID, bidRequest.AuthorType, bidRequest.AuthorID, bidRequest.CreatorUsername, bidRequest.Version, bidRequest.CreatedAt, bidRequest.UpdatedAt,
		).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "status", "tender_id", "author_type", "author_id", "creator_username", "version", "created_at", "updated_at"}).
			AddRow(bidRequest.ID, bidRequest.Name, bidRequest.Description, bidRequest.Status, bidRequest.TenderID, bidRequest.AuthorType, bidRequest.AuthorID, bidRequest.CreatorUsername, bidRequest.Version, bidRequest.CreatedAt, bidRequest.UpdatedAt))
		mock.ExpectCommit()

		bid, err := repo.CreateBid(ctx, bidRequest)
		assert.NoError(t, err)
		assert.NotNil(t, bid)
		assert.Equal(t, bidRequest.ID, bid.ID)
		assert.Equal(t, bidRequest.Name, bid.Name)
		assert.Equal(t, bidRequest.Description, bid.Description)
		assert.Equal(t, bidRequest.Status, bid.Status)
		assert.Equal(t, bidRequest.TenderID, bid.TenderID)
		assert.Equal(t, bidRequest.AuthorType, bid.AuthorType)
		assert.Equal(t, bidRequest.AuthorID, bid.AuthorID)
		assert.Equal(t, bidRequest.CreatorUsername, bid.CreatorUsername)
		assert.Equal(t, bidRequest.Version, bid.Version)
		assert.WithinDuration(t, bidRequest.CreatedAt, bid.CreatedAt, time.Second)
		assert.WithinDuration(t, bidRequest.UpdatedAt, bid.UpdatedAt, time.Second)

		mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidRequest := &model.Bid{
			ID:              uuid.New().String(),
			Name:            "Test Bid",
			Description:     "Test Description",
			Status:          "open",
			TenderID:        uuid.New().String(),
			AuthorType:      "user",
			AuthorID:        uuid.New().String(),
			CreatorUsername: "testuser",
			Version:         1,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`
		INSERT INTO bid (id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
	`)).ExpectQuery().WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		bid, err := repo.CreateBid(ctx, bidRequest)
		assert.Error(t, err)
		assert.Nil(t, bid)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})

}

func TestGetBidById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()

		id := uuid.New().String()
		mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE id = $1`)).ExpectQuery().WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "status", "tender_id", "author_type", "author_id", "creator_username", "version", "created_at", "updated_at",
		}).AddRow(
			id, "Test Bid", "Test Description", "open", uuid.New().String(), "user", uuid.New().String(), "testuser", 1, time.Now(), time.Now(),
		))

		bid, err := repo.GetBidById(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, bid)

		mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()

		id := uuid.New().String()
		mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE id = $1`)).ExpectQuery().WithArgs(id).WillReturnError(sql.ErrConnDone)

		bid, err := repo.GetBidById(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, bid)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()

		id := uuid.New().String()
		mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE id = $1`)).ExpectQuery().WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "status", "tender_id", "author_type", "author_id", "creator_username", "version", "created_at", "updated_at",
		}))

		bid, err := repo.GetBidById(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, bid)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})

}
func TestGetBidByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		username := "testuser"
		limit := 10
		offset := 0

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE creator_username = $1
		LIMIT $2 OFFSET $3
		`)).ExpectQuery().WithArgs(username, limit, offset).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "status", "tender_id", "author_type", "author_id", "creator_username", "version", "created_at", "updated_at",
		}).AddRow(
			uuid.New().String(), "Test Bid", "Test Description", "open", uuid.New().String(), "user", uuid.New().String(), username, 1, time.Now(), time.Now(),
		))

		bids, err := repo.GetBidByUsername(ctx, limit, offset, username)
		assert.NoError(t, err)
		assert.NotNil(t, bids)
		assert.Len(t, bids, 1)

		mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		username := "testuser"
		limit := 10
		offset := 0

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE creator_username = $1
		LIMIT $2 OFFSET $3
		`)).ExpectQuery().WithArgs(username, limit, offset).WillReturnError(sql.ErrConnDone)

		bids, err := repo.GetBidByUsername(ctx, limit, offset, username)
		assert.Error(t, err)
		assert.Nil(t, bids)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		username := "testuser"
		limit := 10
		offset := 0

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE creator_username = $1
		LIMIT $2 OFFSET $3
		`)).ExpectQuery().WithArgs(username, limit, offset).WithArgs(username, limit, offset).WillReturnRows(sqlmock.NewRows([]string{}))

		bids, err := repo.GetBidByUsername(ctx, limit, offset, username)
		assert.Nil(t, bids)
		assert.Len(t, bids, 0)

		mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
func TestGetBidStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidID := uuid.New().String()
		expectedStatus := model.BidStatus("open")

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT status
		FROM bid
		WHERE id = $1
		`)).ExpectQuery().WithArgs(bidID).WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(expectedStatus))

		status, err := repo.GetBidStatus(ctx, bidID)
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, status)

		mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidID := uuid.New().String()

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT status
		FROM bid
		WHERE id = $1
		`)).ExpectQuery().WithArgs(bidID).WillReturnError(sql.ErrConnDone)

		status, err := repo.GetBidStatus(ctx, bidID)
		assert.Error(t, err)
		assert.Equal(t, model.BidStatus(""), status)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidID := uuid.New().String()

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT status
		FROM bid
		WHERE id = $1
		`)).ExpectQuery().WithArgs(bidID).WillReturnRows(sqlmock.NewRows([]string{}))

		status, err := repo.GetBidStatus(ctx, bidID)
		assert.Error(t, err)
		assert.Equal(t, model.BidStatus(""), status)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})
}
func TestUpdateBid(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bid := &model.Bid{
			ID:              uuid.New().String(),
			Name:            "Test Bid",
			Description:     "Test Description",
			Status:          "open",
			TenderID:        uuid.New().String(),
			AuthorType:      "user",
			AuthorID:        uuid.New().String(),
			CreatorUsername: "testuser",
			Version:         1,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		updatedBid := *bid
		updatedBid.Version++
		updatedBid.UpdatedAt = time.Now()

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`
		UPDATE bid 
		SET name = $1, description = $2, status = $3, tender_id = $4, 
			author_type = $5, author_id = $6, creator_username = $7, 
			version = $8, created_at = $9, updated_at = $10 
		WHERE id = $11
		RETURNING id, name, description, status, tender_id, author_type, 
			author_id, creator_username, version, created_at, updated_at
	`)).ExpectQuery().WithArgs(
			bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.CreatorUsername, bid.Version+1, bid.CreatedAt, sqlmock.AnyArg(), bid.ID,
		).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "status", "tender_id", "author_type", "author_id", "creator_username", "version", "created_at", "updated_at"}).
			AddRow(updatedBid.ID, updatedBid.Name, updatedBid.Description, updatedBid.Status, updatedBid.TenderID, updatedBid.AuthorType, updatedBid.AuthorID, updatedBid.CreatorUsername, updatedBid.Version, updatedBid.CreatedAt, updatedBid.UpdatedAt))
		mock.ExpectPrepare(regexp.QuoteMeta(`
		INSERT INTO bid_history (
			id, bid_id, name, description, status, tender_id, 
			author_type, author_id, creator_username, 
			version, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`)).ExpectExec().WithArgs(
			sqlmock.AnyArg(), bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.CreatorUsername, bid.Version, bid.CreatedAt, sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result, err := repo.UpdateBid(ctx, bid)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedBid.ID, result.ID)
		assert.Equal(t, updatedBid.Name, result.Name)
		assert.Equal(t, updatedBid.Description, result.Description)
		assert.Equal(t, updatedBid.Status, result.Status)
		assert.Equal(t, updatedBid.TenderID, result.TenderID)
		assert.Equal(t, updatedBid.AuthorType, result.AuthorType)
		assert.Equal(t, updatedBid.AuthorID, result.AuthorID)
		assert.Equal(t, updatedBid.CreatorUsername, result.CreatorUsername)
		assert.Equal(t, updatedBid.Version, result.Version)
		assert.WithinDuration(t, updatedBid.CreatedAt, result.CreatedAt, time.Second)
		assert.WithinDuration(t, updatedBid.UpdatedAt, result.UpdatedAt, time.Second)

		mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bid := &model.Bid{
			ID:              uuid.New().String(),
			Name:            "Test Bid",
			Description:     "Test Description",
			Status:          "open",
			TenderID:        uuid.New().String(),
			AuthorType:      "user",
			AuthorID:        uuid.New().String(),
			CreatorUsername: "testuser",
			Version:         1,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`
		UPDATE bid 
		SET name = $1, description = $2, status = $3, tender_id = $4, 
			author_type = $5, author_id = $6, creator_username = $7, 
			version = $8, created_at = $9, updated_at = $10 
		WHERE id = $11
		RETURNING id, name, description, status, tender_id, author_type, 
			author_id, creator_username, version, created_at, updated_at
	`)).ExpectQuery().WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		result, err := repo.UpdateBid(ctx, bid)
		assert.Error(t, err)
		assert.Nil(t, result)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})
}
func TestRollbackBidVersion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidID := uuid.New().String()
		version := 1

		historyBid := &model.Bid{
			ID:              uuid.New().String(),
			Name:            "Test Bid",
			Description:     "Test Description",
			Status:          "open",
			TenderID:        uuid.New().String(),
			AuthorType:      "user",
			AuthorID:        uuid.New().String(),
			CreatorUsername: "testuser",
			Version:         version,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, description, status, tender_id, author_type, 
			author_id, creator_username, version, created_at, updated_at
		FROM bid_history
		WHERE bid_id = $1 AND version = $2
		`)).WithArgs(bidID, version).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "status", "tender_id", "author_type", "author_id", "creator_username", "version", "created_at", "updated_at",
		}).AddRow(
			historyBid.ID, historyBid.Name, historyBid.Description, historyBid.Status, historyBid.TenderID, historyBid.AuthorType, historyBid.AuthorID, historyBid.CreatorUsername, historyBid.Version, historyBid.CreatedAt, historyBid.UpdatedAt,
		))

		mock.ExpectPrepare(regexp.QuoteMeta(`
		UPDATE bid
		SET name = $1, description = $2, status = $3, tender_id = $4, 
			author_type = $5, author_id = $6, creator_username = $7, 
			version = $8, created_at = $9, updated_at = $10
		WHERE id = $11
		RETURNING id, name, description, status, tender_id, author_type, 
			author_id, creator_username, version, created_at, updated_at
		`)).ExpectQuery().WithArgs(
			historyBid.Name, historyBid.Description, historyBid.Status, historyBid.TenderID, historyBid.AuthorType, historyBid.AuthorID, historyBid.CreatorUsername, historyBid.Version, historyBid.CreatedAt, sqlmock.AnyArg(), bidID,
		).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "status", "tender_id", "author_type", "author_id", "creator_username", "version", "created_at", "updated_at",
		}).AddRow(
			historyBid.ID, historyBid.Name, historyBid.Description, historyBid.Status, historyBid.TenderID, historyBid.AuthorType, historyBid.AuthorID, historyBid.CreatorUsername, historyBid.Version, historyBid.CreatedAt, time.Now(),
		))

		mock.ExpectCommit()

		updatedBid, err := repo.RollbackBidVersion(ctx, bidID, version)
		assert.NoError(t, err)
		assert.NotNil(t, updatedBid)
		assert.Equal(t, historyBid.ID, updatedBid.ID)
		assert.Equal(t, historyBid.Name, updatedBid.Name)
		assert.Equal(t, historyBid.Description, updatedBid.Description)
		assert.Equal(t, historyBid.Status, updatedBid.Status)
		assert.Equal(t, historyBid.TenderID, updatedBid.TenderID)
		assert.Equal(t, historyBid.AuthorType, updatedBid.AuthorType)
		assert.Equal(t, historyBid.AuthorID, updatedBid.AuthorID)
		assert.Equal(t, historyBid.CreatorUsername, updatedBid.CreatorUsername)
		assert.Equal(t, historyBid.Version, updatedBid.Version)
		assert.WithinDuration(t, historyBid.CreatedAt, updatedBid.CreatedAt, time.Second)

		mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidID := uuid.New().String()
		version := 1

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, description, status, tender_id, author_type, 
			author_id, creator_username, version, created_at, updated_at
		FROM bid_history
		WHERE bid_id = $1 AND version = $2
		`)).WithArgs(bidID, version).WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		updatedBid, err := repo.RollbackBidVersion(ctx, bidID, version)
		assert.Error(t, err)
		assert.Nil(t, updatedBid)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidID := uuid.New().String()
		version := 1

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, description, status, tender_id, author_type, 
			author_id, creator_username, version, created_at, updated_at
		FROM bid_history
		WHERE bid_id = $1 AND version = $2
		`)).WithArgs(bidID, version).WillReturnRows(sqlmock.NewRows([]string{}))
		mock.ExpectRollback()

		updatedBid, err := repo.RollbackBidVersion(ctx, bidID, version)
		assert.Error(t, err)
		assert.Nil(t, updatedBid)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})
}
func TestAddBidFeedback(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidID := uuid.New().String()
		username := "testuser"
		review := "Great bid!"

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`
		INSERT INTO bid_feedback (id, bid_id, username, review, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		`)).ExpectExec().WithArgs(sqlmock.AnyArg(), bidID, username, review, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectPrepare(regexp.QuoteMeta(`
		UPDATE bid
		SET version = version + 1
		WHERE id = $1
		`)).ExpectExec().WithArgs(bidID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, description, status, tender_id, author_type, author_id, creator_username, version, created_at, updated_at
		FROM bid
		WHERE id = $1`)).ExpectQuery().WithArgs(bidID).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "status", "tender_id", "author_type", "author_id", "creator_username", "version", "created_at", "updated_at",
		}).AddRow(
			bidID, "Test Bid", "Test Description", "open", uuid.New().String(), "user", uuid.New().String(), username, 1, time.Now(), time.Now(),
		))

		bid, err := repo.AddBidFeedback(ctx, bidID, username, review)
		assert.NoError(t, err)
		assert.NotNil(t, bid)
		assert.Equal(t, bidID, bid.ID)
		assert.Equal(t, username, bid.CreatorUsername)

		mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		db, mock, repo := setupTestBid(t)
		defer db.Close()

		ctx := context.Background()
		bidID := uuid.New().String()
		username := "testuser"
		review := "Great bid!"

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`
		INSERT INTO bid_feedback (id, bid_id, username, review, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		`)).ExpectExec().WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		bid, err := repo.AddBidFeedback(ctx, bidID, username, review)
		assert.Error(t, err)
		assert.Nil(t, bid)

		mock.ExpectationsWereMet()
		if err == nil {
			t.Errorf("expected error, got nothing")
		}
	})
}
