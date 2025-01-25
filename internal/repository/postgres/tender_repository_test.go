package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTender(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db,mock,err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database", err)
		}
		defer db.Close()

		repo := NewTenderRepository(db, nil)

		tender := &model.Tender{
			ID: "1",
			Name: "test",
			Description: "test",
			ServiceType: "test",
			OrganizationID: "test",
			CreatorUsername: "test",
			Status: "test",
			Version: 1,
		}


		expectedQuery := regexp.QuoteMeta(`
			INSERT INTO tender (id, name, description, service_type, organization_id, creator_username, status, version)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		`)

		rows := sqlmock.NewRows([]string{"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at"}).
			AddRow(tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorUsername, tender.Status, tender.Version, time.Now(), time.Now())

		mock.ExpectPrepare(expectedQuery).ExpectQuery().WithArgs(
			tender.ID,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.OrganizationID,
			tender.CreatorUsername,
			tender.Status,
			tender.Version,
		).WillReturnRows(rows)

		createdTender, err := repo.CreateTender(context.Background(), tender)
		assert.NoError(t, err)
		assert.NotNil(t, createdTender)
		assert.EqualValues(t, tender, createdTender)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
			
	})
	t.Run("error", func(t *testing.T) {
		db,mock,err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database", err)
		}
		defer db.Close()

		repo := NewTenderRepository(db, nil)

		tender := &model.Tender{
			ID: "1",
			Name: "test",
			Description: "test",
			ServiceType: "test",
			OrganizationID: "test",
			CreatorUsername: "test",
			Status: "test",
			Version: 1,
		}


		expectedQuery := regexp.QuoteMeta(`
			INSERT INTO tender (id, name, description, service_type, organization_id, creator_username, status, version)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		`)

		mock.ExpectPrepare(expectedQuery).WillReturnError(errors.New("prepare error"))

		createdTender, err := repo.CreateTender(context.Background(), tender)
		assert.Error(t, err)
		assert.Nil(t, createdTender)
		assert.True(t, strings.Contains(err.Error(), "failed to prepare statement for creating tender"))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}