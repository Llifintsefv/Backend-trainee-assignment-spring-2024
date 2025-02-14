package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func TestGetTenders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tender := model.Tender{}
		ctx := context.Background()
		limit := 10
		offset := 0
		serviceTypes := []model.TenderServiceType{model.TenderServiceTypeConstruction, model.TenderServiceTypeDelivery, model.TenderServiceTypeManufacture}

		expectedQuery := mock.ExpectPrepare(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at FROM tender WHERE service_type = ANY\(\$1\) LIMIT \$2 OFFSET \$3`)

		expectedQuery.ExpectQuery().
			WithArgs(pq.Array(serviceTypes), limit, offset).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
			}).AddRow(
				tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorUsername, tender.Status, tender.Version, time.Now(), time.Now(),
			))

		tenders, err := repo.GetTenders(ctx, limit, offset, serviceTypes)

		assert.NoError(t, err)
		assert.NotNil(t, tenders)

		if len(tenders) == 0 {
			t.Errorf("expected at least one tender, got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("success_empty_result", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()
		ctx := context.Background()
		limit := 10
		offset := 0
		serviceTypes := []model.TenderServiceType{model.TenderServiceTypeConstruction}

		serviceTypeStrings := make([]string, len(serviceTypes))
		for i, st := range serviceTypes {
			serviceTypeStrings[i] = string(st)
		}

		expectedQuery := mock.ExpectPrepare(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at FROM tender WHERE service_type = ANY\(\$1\) LIMIT \$2 OFFSET \$3`)

		expectedQuery.ExpectQuery().
			WithArgs(pq.Array(serviceTypeStrings), limit, offset).
			WillReturnError(sql.ErrNoRows)

		tenders, err := repo.GetTenders(ctx, limit, offset, serviceTypes)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows, but got: %v", err)
		assert.Nil(t, tenders)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("prepare_statement_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		ctx := context.Background()
		limit := 10
		offset := 0
		serviceTypes := []model.TenderServiceType{model.TenderServiceTypeConstruction}

		mock.ExpectPrepare(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at FROM tender WHERE service_type = ANY\(\$1\) LIMIT \$2 OFFSET \$3`).
			WillReturnError(fmt.Errorf("some error"))

		tenders, err := repo.GetTenders(ctx, limit, offset, serviceTypes)

		assert.Error(t, err)
		assert.Nil(t, tenders)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("query_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		ctx := context.Background()
		limit := 10
		offset := 0
		serviceTypes := []model.TenderServiceType{model.TenderServiceTypeConstruction}

		serviceTypeStrings := make([]string, len(serviceTypes))
		for i, st := range serviceTypes {
			serviceTypeStrings[i] = string(st)
		}

		expectedQuery := mock.ExpectPrepare(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at FROM tender WHERE service_type = ANY\(\$1\) LIMIT \$2 OFFSET \$3`)

		expectedQuery.ExpectQuery().
			WithArgs(pq.Array(serviceTypeStrings), limit, offset).
			WillReturnError(fmt.Errorf("some query error"))

		tenders, err := repo.GetTenders(ctx, limit, offset, serviceTypes)

		assert.Error(t, err)
		assert.Nil(t, tenders)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("scan_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		ctx := context.Background()
		limit := 10
		offset := 0
		serviceTypes := []model.TenderServiceType{model.TenderServiceTypeConstruction}

		serviceTypeStrings := make([]string, len(serviceTypes))
		for i, st := range serviceTypes {
			serviceTypeStrings[i] = string(st)
		}

		expectedQuery := mock.ExpectPrepare(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at FROM tender WHERE service_type = ANY\(\$1\) LIMIT \$2 OFFSET \$3`)

		expectedQuery.ExpectQuery().
			WithArgs(pq.Array(serviceTypeStrings), limit, offset).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at",
			}).AddRow( // Missing "updated_at"
				"1", "Test Tender", "Description", "Construction", "123", "user1", "Active", 1, time.Now(),
			))

		tenders, err := repo.GetTenders(ctx, limit, offset, serviceTypes)

		assert.Error(t, err)
		assert.Nil(t, tenders)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("no_service_types_provided", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		ctx := context.Background()
		limit := 10
		offset := 0
		var serviceTypes []model.TenderServiceType

		expectedQuery := mock.ExpectPrepare(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at FROM tender WHERE service_type = ANY\(\$1\) LIMIT \$2 OFFSET \$3`)

		expectedQuery.ExpectQuery().
			WithArgs(pq.Array([]string{}), limit, offset). // Pass an empty string array
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
			}))

		tenders, err := repo.GetTenders(ctx, limit, offset, serviceTypes)

		assert.NoError(t, err)
		assert.Empty(t, tenders) // Expect an empty slice, not nil

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetTenderById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		id := "test"
		ctx := context.Background()

		expectQuery := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		FROM tender
		WHERE id = $1`))

		expectQuery.ExpectQuery().WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
		}).AddRow(
			"id", "Test Tender", "Description", "Construction", "123", "user1", "Active", 1, time.Now(), time.Now()))

		tender, err := repo.GetTenderById(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, tender)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		id := "test"
		ctx := context.Background()

		expectQuery := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		FROM tender
		WHERE id = $1`))

		expectQuery.ExpectQuery().WithArgs(id).WillReturnError(sql.ErrNoRows)
		tender, err := repo.GetTenderById(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, tender)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("prepare_statement_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		id := "test"
		ctx := context.Background()

		expectQuery := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		FROM tender
		WHERE id = $1`))

		expectQuery.ExpectQuery().WithArgs(id).WillReturnError(sql.ErrConnDone)
		tender, err := repo.GetTenderById(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, tender)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
func TestGetTenderByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		username := "testuser"
		limit := 10
		offset := 0
		ctx := context.Background()

		expectedQuery := mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
			FROM tender
			WHERE creator_username = $1
			LIMIT $2 OFFSET $3
		`))

		expectedQuery.ExpectQuery().
			WithArgs(username, limit, offset).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
			}).AddRow(
				"1", "Test Tender", "Description", "ServiceType", "OrgID", username, "Status", 1, time.Now(), time.Now(),
			))

		tenders, err := repo.GetTenderByUsername(ctx, limit, offset, username)

		assert.NoError(t, err)
		assert.NotNil(t, tenders)
		assert.Len(t, tenders, 1)
		assert.Equal(t, username, tenders[0].CreatorUsername)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("prepare_statement_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		username := "testuser"
		limit := 10
		offset := 0
		ctx := context.Background()

		mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
			FROM tender
			WHERE creator_username = $1
			LIMIT $2 OFFSET $3
		`)).WillReturnError(fmt.Errorf("prepare statement error"))

		tenders, err := repo.GetTenderByUsername(ctx, limit, offset, username)

		assert.Error(t, err)
		assert.Nil(t, tenders)
		assert.EqualError(t, err, "failed to prepare statement for getting tenders: prepare statement error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("query_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		username := "testuser"
		limit := 10
		offset := 0
		ctx := context.Background()

		expectedQuery := mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
			FROM tender
			WHERE creator_username = $1
			LIMIT $2 OFFSET $3
		`))

		expectedQuery.ExpectQuery().
			WithArgs(username, limit, offset).
			WillReturnError(fmt.Errorf("query error"))

		tenders, err := repo.GetTenderByUsername(ctx, limit, offset, username)

		assert.Error(t, err)
		assert.Nil(t, tenders)
		assert.EqualError(t, err, "failed to execute query for getting tenders: query error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("scan_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		username := "testuser"
		limit := 10
		offset := 0
		ctx := context.Background()

		expectedQuery := mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
			FROM tender
			WHERE creator_username = $1
			LIMIT $2 OFFSET $3
		`))

		expectedQuery.ExpectQuery().
			WithArgs(username, limit, offset).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at",
			}).AddRow( // Missing "updated_at"
				"1", "Test Tender", "Description", "ServiceType", "OrgID", username, "Status", 1, time.Now(),
			))

		tenders, err := repo.GetTenderByUsername(ctx, limit, offset, username)

		assert.Error(t, err)
		assert.Nil(t, tenders)
		assert.Contains(t, err.Error(), "failed to scan tender")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
func TestIsUserResponsibleForTender(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tenderID := "test-tender-id"
		username := "testuser"
		ctx := context.Background()

		expectedQuery := mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT EXISTS (
				SELECT 1
				FROM organization_responsible orr
				JOIN tender t ON t.organization_id = orr.organization_id
				JOIN employee e ON e.id = orr.user_id
				WHERE t.id = $1 AND e.username = $2
			)
		`))

		expectedQuery.ExpectQuery().
			WithArgs(tenderID, username).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		isResponsible, err := repo.IsUserResponsibleForTender(ctx, tenderID, username)

		assert.NoError(t, err)
		assert.True(t, isResponsible)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("not_responsible", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tenderID := "test-tender-id"
		username := "testuser"
		ctx := context.Background()

		expectedQuery := mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT EXISTS (
				SELECT 1
				FROM organization_responsible orr
				JOIN tender t ON t.organization_id = orr.organization_id
				JOIN employee e ON e.id = orr.user_id
				WHERE t.id = $1 AND e.username = $2
			)
		`))

		expectedQuery.ExpectQuery().
			WithArgs(tenderID, username).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		isResponsible, err := repo.IsUserResponsibleForTender(ctx, tenderID, username)

		assert.NoError(t, err)
		assert.False(t, isResponsible)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("prepare_statement_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tenderID := "test-tender-id"
		username := "testuser"
		ctx := context.Background()

		mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT EXISTS (
				SELECT 1
				FROM organization_responsible orr
				JOIN tender t ON t.organization_id = orr.organization_id
				JOIN employee e ON e.id = orr.user_id
				WHERE t.id = $1 AND e.username = $2
			)
		`)).WillReturnError(fmt.Errorf("prepare statement error"))

		isResponsible, err := repo.IsUserResponsibleForTender(ctx, tenderID, username)

		assert.Error(t, err)
		assert.False(t, isResponsible)
		assert.EqualError(t, err, "error preparing statement for checking user is responsible: prepare statement error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("query_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tenderID := "test-tender-id"
		username := "testuser"
		ctx := context.Background()

		expectedQuery := mock.ExpectPrepare(regexp.QuoteMeta(`
			SELECT EXISTS (
				SELECT 1
				FROM organization_responsible orr
				JOIN tender t ON t.organization_id = orr.organization_id
				JOIN employee e ON e.id = orr.user_id
				WHERE t.id = $1 AND e.username = $2
			)
		`))

		expectedQuery.ExpectQuery().
			WithArgs(tenderID, username).
			WillReturnError(fmt.Errorf("query error"))

		isResponsible, err := repo.IsUserResponsibleForTender(ctx, tenderID, username)

		assert.Error(t, err)
		assert.False(t, isResponsible)
		assert.EqualError(t, err, "error checking user is responsible: query error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

}
func TestUpdateTender(t *testing.T) {
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
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mock.ExpectBegin()

		mock.ExpectPrepare(regexp.QuoteMeta(`
			UPDATE tender
			SET name = $2, description = $3, service_type = $4, organization_id = $5, creator_username = $6, status = $7, version = $8, updated_at = $9
			WHERE id = $1
			RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		`)).ExpectQuery().WithArgs(
			tender.ID,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.OrganizationID,
			tender.CreatorUsername,
			tender.Status,
			tender.Version+1,
			sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
		}).AddRow(
			tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorUsername, tender.Status, tender.Version+1, tender.CreatedAt, time.Now(),
		))

		mock.ExpectPrepare(regexp.QuoteMeta(`
			INSERT INTO tender_history (id, tender_id, name, description, service_type, status, organization_id, creator_username, version, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`)).ExpectExec().WithArgs(
			sqlmock.AnyArg(),
			tender.ID,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.Status,
			tender.OrganizationID,
			tender.CreatorUsername,
			tender.Version,
			tender.CreatedAt,
			sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		updatedTender, err := repo.UpdateTender(context.Background(), &tender)

		assert.NoError(t, err)
		assert.NotNil(t, updatedTender)
		assert.Equal(t, tender.Name, updatedTender.Name)
		assert.Equal(t, tender.Version+1, updatedTender.Version)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("begin_transaction_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tender := &model.Tender{ID: uuid.New().String(), Name: "Test Tender"}

		mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))

		_, err := repo.UpdateTender(context.Background(), tender)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to begin transaction: begin transaction error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("prepare_statement_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tender := &model.Tender{ID: uuid.New().String(), Name: "Test Tender"}

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE tender`)).WillReturnError(errors.New("prepare statement error"))
		mock.ExpectRollback()

		_, err := repo.UpdateTender(context.Background(), tender)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to prepare statement for updating tender: prepare statement error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("query_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tender := &model.Tender{ID: uuid.New().String(), Name: "Test Tender"}

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE tender`)).
			ExpectQuery().
			WillReturnError(errors.New("query error"))
		mock.ExpectRollback()

		_, err := repo.UpdateTender(context.Background(), tender)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to scan updated tender: query error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("insert_history_error", func(t *testing.T) {
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
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mock.ExpectBegin()

		mock.ExpectPrepare(regexp.QuoteMeta(`
			UPDATE tender
			SET name = $2, description = $3, service_type = $4, organization_id = $5, creator_username = $6, status = $7, version = $8, updated_at = $9
			WHERE id = $1
			RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		`)).ExpectQuery().WithArgs(
			tender.ID,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.OrganizationID,
			tender.CreatorUsername,
			tender.Status,
			tender.Version+1,
			sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
		}).AddRow(
			tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorUsername, tender.Status, tender.Version+1, tender.CreatedAt, time.Now(),
		))

		mock.ExpectPrepare(regexp.QuoteMeta(`
			INSERT INTO tender_history (id, tender_id, name, description, service_type, status, organization_id, creator_username, version, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`)).ExpectExec().WithArgs(
			sqlmock.AnyArg(),
			tender.ID,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.Status,
			tender.OrganizationID,
			tender.CreatorUsername,
			tender.Version,
			tender.CreatedAt,
			sqlmock.AnyArg(),
		).WillReturnError(errors.New("insert history error"))
		mock.ExpectRollback()

		_, err := repo.UpdateTender(context.Background(), &tender)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to insert tender history: insert history error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("commit_error", func(t *testing.T) {
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
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mock.ExpectBegin()

		mock.ExpectPrepare(regexp.QuoteMeta(`
			UPDATE tender
			SET name = $2, description = $3, service_type = $4, organization_id = $5, creator_username = $6, status = $7, version = $8, updated_at = $9
			WHERE id = $1
			RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		`)).ExpectQuery().WithArgs(
			tender.ID,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.OrganizationID,
			tender.CreatorUsername,
			tender.Status,
			tender.Version+1,
			sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
		}).AddRow(
			tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorUsername, tender.Status, tender.Version+1, tender.CreatedAt, time.Now(),
		))

		mock.ExpectPrepare(regexp.QuoteMeta(`
			INSERT INTO tender_history (id, tender_id, name, description, service_type, status, organization_id, creator_username, version, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`)).ExpectExec().WithArgs(
			sqlmock.AnyArg(),
			tender.ID,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.Status,
			tender.OrganizationID,
			tender.CreatorUsername,
			tender.Version,
			tender.CreatedAt,
			sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		_, err := repo.UpdateTender(context.Background(), &tender)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to commit transaction: commit error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
func TestRollbackTenderVersion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tenderID := "test-tender-id"
		version := 1
		ctx := context.Background()

		historyTender := model.Tender{
			ID:              tenderID,
			Name:            "Test Tender",
			Description:     "Test Description",
			ServiceType:     "Test Service",
			OrganizationID:  "test-org-id",
			CreatorUsername: "testuser",
			Status:          "test",
			Version:         version,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT tender_id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
			FROM tender_history
			WHERE tender_id = $1 AND version = $2
		`)).WithArgs(tenderID, version).
			WillReturnRows(sqlmock.NewRows([]string{
				"tender_id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
			}).AddRow(
				historyTender.ID, historyTender.Name, historyTender.Description, historyTender.ServiceType, historyTender.OrganizationID, historyTender.CreatorUsername, historyTender.Status, historyTender.Version, historyTender.CreatedAt, historyTender.UpdatedAt,
			))

		mock.ExpectPrepare(regexp.QuoteMeta(`
			UPDATE tender
			SET name = $2, description = $3, service_type = $4, organization_id = $5, creator_username = $6, status = $7, version = $8, updated_at = $9
			WHERE id = $1
			RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
		`)).ExpectQuery().WithArgs(
			historyTender.ID,
			historyTender.Name,
			historyTender.Description,
			historyTender.ServiceType,
			historyTender.OrganizationID,
			historyTender.CreatorUsername,
			historyTender.Status,
			historyTender.Version,
			sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "service_type", "organization_id", "creator_username", "status", "version", "created_at", "updated_at",
		}).AddRow(
			historyTender.ID, historyTender.Name, historyTender.Description, historyTender.ServiceType, historyTender.OrganizationID, historyTender.CreatorUsername, historyTender.Status, historyTender.Version, historyTender.CreatedAt, time.Now(),
		))

		mock.ExpectCommit()

		updatedTender, err := repo.RollbackTenderVersion(ctx, tenderID, version)

		assert.NoError(t, err)
		assert.NotNil(t, updatedTender)
		assert.Equal(t, historyTender.Name, updatedTender.Name)
		assert.Equal(t, historyTender.Version, updatedTender.Version)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("begin_transaction_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tenderID := "test-tender-id"
		version := 1
		ctx := context.Background()

		mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))

		_, err := repo.RollbackTenderVersion(ctx, tenderID, version)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to begin transaction: begin transaction error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("no_rows_error", func(t *testing.T) {
		db, mock, repo := setupTest(t)
		defer db.Close()

		tenderID := "test-tender-id"
		version := 1
		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT tender_id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
			FROM tender_history
			WHERE tender_id = $1 AND version = $2
		`)).WithArgs(tenderID, version).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		_, err := repo.RollbackTenderVersion(ctx, tenderID, version)
		assert.Error(t, err)
		assert.EqualError(t, err, model.ErrVersionNotFound.Error())

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
