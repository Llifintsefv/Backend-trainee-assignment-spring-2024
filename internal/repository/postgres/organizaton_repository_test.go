package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"log/slog"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestOrganization(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository.OrganizationRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	logger := slog.Default()
	repo := NewOrganizationRepository(db, logger)
	return db, mock, repo
}

func TestGetOrganizationById(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTestOrganization(t)
		defer db.Close()

		ctx := context.Background()
		id := uuid.New().String()

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, description,type,created_at, updated_at
		FROM organization
		WHERE id = $1
		`)).ExpectQuery().WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "type", "created_at", "updated_at",
		}).AddRow(
			id, "Test Organization", "Test Description", "Public", time.Now(), time.Now(),
		))

		organization, err := repo.GetOrganizationById(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, organization)

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, repo := setupTestOrganization(t)
		defer db.Close()

		ctx := context.Background()
		id := uuid.New().String()

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, description,type,created_at, updated_at
		FROM organization
		WHERE id = $1
		`)).ExpectQuery().WithArgs(id).WillReturnError(sql.ErrNoRows)

		organization, err := repo.GetOrganizationById(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, organization)
	})

	t.Run("error", func(t *testing.T) {
		db, mock, repo := setupTestOrganization(t)
		defer db.Close()

		ctx := context.Background()
		id := uuid.New().String()

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT id, name, description,type,created_at, updated_at
		FROM organization
		WHERE id = $1
		`)).ExpectQuery().WithArgs(id).WillReturnError(sql.ErrConnDone)

		organization, err := repo.GetOrganizationById(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, organization)
	})

}
func TestIsUserResponsibleForOrganization(t *testing.T) {

	t.Run("user is responsible", func(t *testing.T) {
		db, mock, repo := setupTestOrganization(t)
		defer db.Close()

		ctx := context.Background()
		organizationID := uuid.New().String()
		username := "testuser"

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT EXISTS (
			SELECT 1
			FROM organization_responsible orr
			JOIN employee e ON e.id = orr.user_id
			WHERE orr.organization_id = $1 AND e.username = $2
		)
		`)).ExpectQuery().WithArgs(organizationID, username).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		isResponsible, err := repo.IsUserResponsibleForOrganization(ctx, organizationID, username)

		assert.NoError(t, err)
		assert.True(t, isResponsible)

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("user is not responsible", func(t *testing.T) {
		db, mock, repo := setupTestOrganization(t)
		defer db.Close()

		ctx := context.Background()
		organizationID := uuid.New().String()
		username := "testuser"

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT EXISTS (
			SELECT 1
			FROM organization_responsible orr
			JOIN employee e ON e.id = orr.user_id
			WHERE orr.organization_id = $1 AND e.username = $2
		)
		`)).ExpectQuery().WithArgs(organizationID, username).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		isResponsible, err := repo.IsUserResponsibleForOrganization(ctx, organizationID, username)

		assert.NoError(t, err)
		assert.False(t, isResponsible)

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, repo := setupTestOrganization(t)
		defer db.Close()

		ctx := context.Background()
		organizationID := uuid.New().String()
		username := "testuser"

		mock.ExpectPrepare(regexp.QuoteMeta(`
		SELECT EXISTS (
			SELECT 1
			FROM organization_responsible orr
			JOIN employee e ON e.id = orr.user_id
			WHERE orr.organization_id = $1 AND e.username = $2
		)
		`)).ExpectQuery().WithArgs(organizationID, username).WillReturnError(sql.ErrConnDone)

		isResponsible, err := repo.IsUserResponsibleForOrganization(ctx, organizationID, username)

		assert.Error(t, err)
		assert.False(t, isResponsible)
	})
}
