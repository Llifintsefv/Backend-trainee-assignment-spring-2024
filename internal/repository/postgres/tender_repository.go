package postgres

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"Backend-trainee-assignment-autumn-2024/internal/repository"
	"context"
	"database/sql"
	"fmt"
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

// func (r *tenderRepository) CreateTender(ctx context.Context, tender *model.Tender) (*model.Tender, error) {
// 	stmt, err := r.db.PrepareContext(ctx, `
// 		INSERT INTO tender (id, name, description, service_type, organization_id, creator_username, status, version)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
// 		RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
// 	`)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to prepare statement for creating tender: %w", err)
// 	}
// 	defer stmt.Close()

// 	row := stmt.QueryRowContext(ctx,
// 		tender.ID,
// 		tender.Name,
// 		tender.Description,
// 		tender.ServiceType,
// 		tender.OrganizationID,
// 		tender.CreatorUsername,
// 		tender.Status,
// 		tender.Version,
// 	)

// 	if err := row.Scan(
// 		&tender.ID,
// 		&tender.Name,
// 		&tender.Description,
// 		&tender.ServiceType,
// 		&tender.OrganizationID,
// 		&tender.CreatorUsername,
// 		&tender.Status,
// 		&tender.Version,
// 		&tender.CreatedAt,
// 		&tender.UpdatedAt,
// 	); err != nil {
// 		r.logger.ErrorContext(ctx, "Error creating tender and scanning result", slog.Any("error", err))
// 		return nil, fmt.Errorf("failed to execute query and scan result: %w", err)
// 	}

// 	return tender, nil
// }

func (r *tenderRepository) CreateTender(ctx context.Context, tender *model.Tender) (*model.Tender, error) {
	r.logger.DebugContext(ctx, "--- НАЧАЛО ВЫВОДА ДАННЫХ ИЗ БАЗЫ ДАННЫХ (ОТЛАДКА) ---")

	tablesToDebug := []string{"organization", "tender"} // Список таблиц для вывода

	for _, tableName := range tablesToDebug {
		r.logger.DebugContext(ctx, "--- Содержимое таблицы:", slog.String("table", tableName), slog.String("---", "---"))
		rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %s", tableName))
		if err != nil {
			r.logger.ErrorContext(ctx, "Ошибка запроса к таблице", slog.String("table", tableName), slog.Any("error", err))
			continue // Переходим к следующей таблице в случае ошибки
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			r.logger.ErrorContext(ctx, "Ошибка получения списка колонок таблицы", slog.String("table", tableName), slog.Any("error", err))
			continue
		}

		for rows.Next() {
			columnValues := make([]interface{}, len(columns))
			columnPointers := make([]interface{}, len(columns))
			for i := range columns {
				columnPointers[i] = &columnValues[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				r.logger.ErrorContext(ctx, "Ошибка сканирования строки таблицы", slog.String("table", tableName), slog.Any("error", err))
				break // Прерываем вывод текущей таблицы в случае ошибки сканирования
			}

			rowMap := make(map[string]interface{})
			for i, colName := range columns {
				val := columnValues[i]
				// Преобразуем []byte в string для текстовых полей (если нужно)
				if byteSlice, ok := val.([]byte); ok {
					val = string(byteSlice)
				}
				rowMap[colName] = val
			}
			r.logger.DebugContext(ctx, "Строка:", slog.Any("data", rowMap))
		}

		if err := rows.Err(); err != nil {
			r.logger.ErrorContext(ctx, "Ошибка итерации по строкам таблицы", slog.String("table", tableName), slog.Any("error", err))
		}
	}

	r.logger.DebugContext(ctx, "--- КОНЕЦ ВЫВОДА ДАННЫХ ИЗ БАЗЫ ДАННЫХ (ОТЛАДКА) ---")


	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO tender (id, name, description, service_type, organization_id, creator_username, status, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, description, service_type, organization_id, creator_username, status, version, created_at, updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for creating tender: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx,
		tender.ID,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.OrganizationID,
		tender.CreatorUsername,
		tender.Status,
		tender.Version,
	)
	if err := row.Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.OrganizationID,
		&tender.CreatorUsername,
		&tender.Status,
		&tender.Version,
		&tender.CreatedAt,
		&tender.UpdatedAt,
	); err != nil {
		r.logger.ErrorContext(ctx, "Error creating tender and scanning result", slog.Any("error", err))
		return nil, fmt.Errorf("failed to execute query and scan result: %w", err)
	}

	return tender, nil

	// ... (остальная часть вашей функции)
}