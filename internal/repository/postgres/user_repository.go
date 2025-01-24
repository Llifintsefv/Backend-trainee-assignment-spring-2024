package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &organizationRepository{
		db: db,
	}
}