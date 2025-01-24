package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"database/sql"
	"log/slog"
)

type userRepository struct {
	db *sql.DB
	logger *slog.Logger
}

func NewUserRepository(db *sql.DB, logger *slog.Logger) repository.UserRepository {
	return &userRepository{
		db: db,
		logger: logger,
	}
}