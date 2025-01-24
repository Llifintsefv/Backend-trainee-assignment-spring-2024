package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"database/sql"
)

type tenderRepository struct {
	db *sql.DB
}


func NewTenderRepository(db *sql.DB) repository.TenderRepository {
	return &tenderRepository{
		db: db,
	}
}

