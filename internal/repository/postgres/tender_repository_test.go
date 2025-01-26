package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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


func TestGetTenders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewTenderRepository(db, nil)
		ctx := context.Background()

		mockRows := sqlmock.NewRows([]string{
			"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
		}).AddRow(
			"1", "Test Tender", "Description", "Construction", "Org1", "User1", "CREATED", 1, time.Now(), time.Now(),
		).AddRow(
			"2", "Another Tender", "Description", "Delivery", "Org2", "User2", "PUBLISHED", 1, time.Now(), time.Now(),
		)

		mock.ExpectPrepare("^SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at FROM tender WHERE service_type = ANY\\(\\$1\\) LIMIT \\$2 OFFSET \\$3$").
			ExpectQuery().
			WithArgs(sqlmock.AnyArg(), 10, 0).
			WillReturnRows(mockRows)


		
		serviceTypes := []model.TenderServiceType{model.TenderServiceTypeConstruction, model.TenderServiceTypeDelivery}
		tenders, err := repo.GetTenders(ctx, 10, 0, serviceTypes)

		
		require.NoError(t, err)
		require.Len(t, tenders, 2)
		assert.Equal(t, "1", tenders[0].ID)
		assert.Equal(t, "Test Tender", tenders[0].Name)
		assert.Equal(t, model.TenderServiceTypeConstruction, tenders[0].ServiceType)
		assert.Equal(t, "Org1", tenders[0].OrganizationID)
	})
	t.Run("returns error on query preparation", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		ctx := context.Background()
		mock.ExpectPrepare("^SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at FROM tender WHERE service_type = ANY\\(\\$1\\) LIMIT \\$2 OFFSET \\$3$").
			WillReturnError(sql.ErrConnDone)

		repo := NewTenderRepository(db,nil)

		serviceTypes := []model.TenderServiceType{model.TenderServiceTypeConstruction}
		_, err = repo.GetTenders(ctx, 10, 0, serviceTypes)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to prepare statement for getting tenders")
	})

}
