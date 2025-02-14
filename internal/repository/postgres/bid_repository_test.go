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
	})

}
