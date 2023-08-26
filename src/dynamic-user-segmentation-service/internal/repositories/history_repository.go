package repositories

import (
	"database/sql"
	"dynamic-user-segmentation-service/internal/db"
	"dynamic-user-segmentation-service/internal/errors"
	"time"
)

type HistoryRepository interface {
	SetAddingHistoryRecord(userId int, slug string) error
	SetRemovingHistoryRecord(userId int, slug string) error
	// GetHistoryByDate() ([]*models.Segment, error) TODO: сделать получение истории
}

type PostgresHistoryRepository struct {
	db *sql.DB
}

const (
	saveRecord = `INSERT INTO history (user_id, slug, action_date, operation_type) VALUES ($1, $2, $3, $4);`
)

func (r *PostgresHistoryRepository) SetAddingHistoryRecord(userId int, slug string) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Exec(saveRecord, userId, slug, time.Now(), "ADDING")
	if err != nil {
		return errors.DatabaseWritingError
	}

	return nil
}

func (r *PostgresHistoryRepository) SetRemovingHistoryRecord(userId int, slug string) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Exec(saveRecord, userId, slug, time.Now(), "REMOVING")
	if err != nil {
		return errors.DatabaseWritingError
	}

	return nil
}
