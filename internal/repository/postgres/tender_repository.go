package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"database/sql"
	"log/slog"
)

type tenderRepository struct {
	db *sql.DB
	logger *slog.Logger
}


func NewTenderRepository(db *sql.DB, logger *slog.Logger) repository.TenderRepository {
	return &tenderRepository{
		db: db,
		logger: logger,
	}
}

