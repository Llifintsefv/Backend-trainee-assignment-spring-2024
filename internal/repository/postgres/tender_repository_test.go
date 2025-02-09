package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository.TenderRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	logger := slog.Default()
	repo := NewTenderRepository(db, logger)

	return db, mock, repo
}

func TestCreateTender(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTest(t)

		defer db.Close()

		tender := model.Tender{
			ID:              uuid.New().String(),
			Name:            "Test Tender",
			Description:     "Test Description",
			ServiceType:     "Test Service",
			OrganizationID:  uuid.New().String(),
			CreatorUsername: "testuser",
			Status:          "test",
			Version:         1,
		}

		mock.ExpectBegin()

		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO tender (id, name, description, service_type, organization_id, creator_username, status, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at`)).
			ExpectQuery().WithArgs(
			tender.ID,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.OrganizationID,
			tender.CreatorUsername,
			tender.Status,
			tender.Version,
		).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at"}).
			AddRow(tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorUsername, tender.Status, tender.Version, time.Now(), time.Now()))

		mock.ExpectCommit()

		createdTender, err := repo.CreateTender(context.Background(), &tender)

		assert.NoError(t, err)
		assert.NotNil(t, createdTender)
		assert.Equal(t, tender.Name, createdTender.Name)
		assert.NotEmpty(t, createdTender.UpdatedAt)
		assert.NotEmpty(t, createdTender.CreatedAt)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

	t.Run("data base error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()
		tender := &model.Tender{ID: uuid.New().String(), Name: "Test Tender"}

		mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))

		_, err := repo.CreateTender(context.Background(), tender)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to begin transaction: begin transaction error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("prepare statement error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()
		tender := &model.Tender{ID: uuid.New().String(), Name: "Test Tender"}

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO tender`)).WillReturnError(errors.New("prepare statement error"))
		mock.ExpectRollback()

		_, err := repo.CreateTender(context.Background(), tender)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to prepare statement for creating tender: prepare statement error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tender := &model.Tender{ID: uuid.New().String(), Name: "Test Tender"}

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO tender`)).
			ExpectQuery().
			WillReturnError(errors.New("query error"))
		mock.ExpectRollback()

		_, err := repo.CreateTender(context.Background(), tender)
		assert.Error(t, err)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

}
