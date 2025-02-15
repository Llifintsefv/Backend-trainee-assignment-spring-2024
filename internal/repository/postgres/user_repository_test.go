package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestUser(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repository.UserRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	logger := slog.Default()
	repo := NewUserRepository(db, logger)

	return db, mock, repo
}

func TestGetUserById(t *testing.T) {
	db, mock, repo := setupTestUser(t)
	defer db.Close()

	ctx := context.Background()
	userID := uuid.New().String()
	expectedUser := &model.User{
		Id:         userID,
		Username:   "testuser",
		First_name: "Test",
		Last_name:  "User",
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	query := regexp.QuoteMeta(`
		SELECT id, username, first_name, last_name, created_at, updated_at
		FROM employee  
		WHERE id = $1
	`)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "first_name", "last_name", "created_at", "updated_at"}).
			AddRow(expectedUser.Id, expectedUser.Username, expectedUser.First_name, expectedUser.Last_name, expectedUser.Created_at, expectedUser.Updated_at)

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(userID).WillReturnRows(rows)

		user, err := repo.GetUserById(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectPrepare(query).ExpectQuery().WithArgs(userID).WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserById(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("prepare statement error", func(t *testing.T) {
		mock.ExpectPrepare(query).WillReturnError(fmt.Errorf("prepare error"))

		user, err := repo.GetUserById(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to prepare statement for getting user by id")
	})

	t.Run("query execution error", func(t *testing.T) {
		mock.ExpectPrepare(query).ExpectQuery().WithArgs(userID).WillReturnError(fmt.Errorf("query error"))

		user, err := repo.GetUserById(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to execute query for getting user by id")
	})
}
func TestGetUserByUsername(t *testing.T) {
	db, mock, repo := setupTestUser(t)
	defer db.Close()

	ctx := context.Background()
	username := "testuser"
	expectedUser := &model.User{
		Id:         uuid.New().String(),
		Username:   username,
		First_name: "Test",
		Last_name:  "User",
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	query := regexp.QuoteMeta(`
		SELECT id, username, first_name, last_name, created_at, updated_at
		FROM employee  
		WHERE username = $1
	`)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "first_name", "last_name", "created_at", "updated_at"}).
			AddRow(expectedUser.Id, expectedUser.Username, expectedUser.First_name, expectedUser.Last_name, expectedUser.Created_at, expectedUser.Updated_at)

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(username).WillReturnRows(rows)

		user, err := repo.GetUserByUsername(ctx, username)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectPrepare(query).ExpectQuery().WithArgs(username).WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByUsername(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("prepare statement error", func(t *testing.T) {
		mock.ExpectPrepare(query).WillReturnError(fmt.Errorf("prepare error"))

		user, err := repo.GetUserByUsername(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to prepare statement for getting user by username")
	})

	t.Run("query execution error", func(t *testing.T) {
		mock.ExpectPrepare(query).ExpectQuery().WithArgs(username).WillReturnError(fmt.Errorf("query error"))

		user, err := repo.GetUserByUsername(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to execute query for getting user by username")
	})
}
func TestGetOrganizationByUsername(t *testing.T) {
	db, mock, repo := setupTestUser(t)
	defer db.Close()

	ctx := context.Background()
	username := "testuser"
	expectedOrganization := &model.Organization{
		Id:         uuid.New().String(),
		Name:       "Test Organization",
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	query := regexp.QuoteMeta(`
		SELECT id, name, created_at, updated_at
		FROM organization  
		WHERE username = $1
	`)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(expectedOrganization.Id, expectedOrganization.Name, expectedOrganization.Created_at, expectedOrganization.Updated_at)

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(username).WillReturnRows(rows)

		organization, err := repo.GetOrganizationByUsername(ctx, username)
		assert.NoError(t, err)
		assert.NotNil(t, organization)
		assert.Equal(t, expectedOrganization, organization)
	})

	t.Run("organization not found", func(t *testing.T) {
		mock.ExpectPrepare(query).ExpectQuery().WithArgs(username).WillReturnError(sql.ErrNoRows)

		organization, err := repo.GetOrganizationByUsername(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, organization)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("prepare statement error", func(t *testing.T) {
		mock.ExpectPrepare(query).WillReturnError(fmt.Errorf("prepare error"))

		organization, err := repo.GetOrganizationByUsername(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, organization)
		assert.Contains(t, err.Error(), "failed to prepare statement for getting organization by username")
	})

	t.Run("query execution error", func(t *testing.T) {
		mock.ExpectPrepare(query).ExpectQuery().WithArgs(username).WillReturnError(fmt.Errorf("query error"))

		organization, err := repo.GetOrganizationByUsername(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, organization)
		assert.Contains(t, err.Error(), "failed to execute query for getting organization by username")
	})
}
