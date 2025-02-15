package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type userRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserRepository(db *sql.DB, logger *slog.Logger) repository.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, username, first_name, last_name, created_at, updated_at
		FROM employee  
		WHERE id = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for getting user by id: %w", err)
	}
	defer stmt.Close()

	user := model.User{}
	err = stmt.QueryRowContext(ctx, id).Scan(
		&user.Id,
		&user.Username,
		&user.First_name,
		&user.Last_name,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting user by id", slog.Any("error", err))
			return nil, fmt.Errorf("failed to execute query for getting user by id: %w", err)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, username, first_name, last_name, created_at, updated_at
		FROM employee  
		WHERE username = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for getting user by username: %w", err)
	}
	defer stmt.Close()

	user := model.User{}
	err = stmt.QueryRowContext(ctx, username).Scan(
		&user.Id,
		&user.Username,
		&user.First_name,
		&user.Last_name,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting user by username", slog.Any("error", err))
			return nil, fmt.Errorf("failed to execute query for getting user by username: %w", err)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetOrganizationByUsername(ctx context.Context, username string) (*model.Organization, error) {

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, name, created_at, updated_at
		FROM organization  
		WHERE username = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for getting organization by username: %w", err)
	}
	defer stmt.Close()

	organization := model.Organization{}
	err = stmt.QueryRowContext(ctx, username).Scan(
		&organization.Id,
		&organization.Name,
		&organization.Created_at,
		&organization.Updated_at,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.ErrorContext(ctx, "Error getting organization by username", slog.Any("error", err))
			return nil, fmt.Errorf("failed to execute query for getting organization by username: %w", err)
		}
		return nil, err
	}

	return &organization, nil
}
