package repositories

import (
	goErrors "errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type PostgresHistoryRepository struct {
	db *sqlx.DB
}

var (
	ErrRecordNotFound       = goErrors.New("Record was not found")
	ErrDatabaseWritingError = goErrors.New("Error while writing to DB")
	ErrDatabaseReadingError = goErrors.New("Error while reading from DB")
	ErrRecordAlreadyExists  = goErrors.New("Record with this data already exists")
)

func NewHistoryRepository(db *sqlx.DB) *PostgresHistoryRepository {
	return &PostgresHistoryRepository{
		db: db,
	}
}

const (
	saveRecord = `INSERT INTO history (user_id, slug, action_date, operation_type) VALUES ($1, $2, $3, $4);`
)

func (r *PostgresHistoryRepository) SetAddingHistoryRecord(userId int, slug string) error {
	_, err := r.db.Exec(saveRecord, userId, slug, time.Now(), "ADDING")
	if err != nil {
		return ErrDatabaseWritingError
	}

	return nil
}

func (r *PostgresHistoryRepository) SetRemovingHistoryRecord(userId int, slug string) error {
	_, err := r.db.Exec(saveRecord, userId, slug, time.Now(), "REMOVING")
	if err != nil {
		return ErrDatabaseWritingError
	}

	return nil
}
