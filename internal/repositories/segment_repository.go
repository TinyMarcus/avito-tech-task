package repositories

import (
	"database/sql"
	goErrors "errors"

	"github.com/jmoiron/sqlx"

	"github.com/TinyMarcus/avito-tech-task/internal/models"
)

type PostgresSegmentRepository struct {
	db *sqlx.DB
}

func NewSegmentRepository(db *sqlx.DB) *PostgresSegmentRepository {
	return &PostgresSegmentRepository{
		db: db,
	}
}

const (
	selectSegments       = `SELECT id, slug, description FROM segments;`
	selectSegmentBySlug  = `SELECT id, slug, description FROM segments WHERE Slug = $1;`
	createSegment        = `INSERT INTO segments (slug, description) VALUES ($1, $2) RETURNING slug;`
	updateSegment        = `UPDATE segments SET description = $1 WHERE slug = $2;`
	deleteSegment        = `DELETE FROM segments WHERE slug = $1 RETURNING slug;`
	checkIfSegmentExists = `SELECT id, slug, description FROM segments WHERE slug = $1;`
)

func (r *PostgresSegmentRepository) GetAllSegments() ([]*models.Segment, error) {
	var segments []*models.Segment

	rows, err := r.db.Query(selectSegments)
	if err != nil {
		return nil, ErrDatabaseReadingError
	}

	if err := rows.Err(); err != nil {
		return nil, ErrDatabaseReadingError
	}

	for rows.Next() {
		segment := new(models.Segment)
		if err := rows.Scan(&segment.Id, &segment.Slug, &segment.Description); err != nil {
			return nil, ErrDatabaseReadingError
		}
		segments = append(segments, segment)
	}

	defer rows.Close()
	return segments, nil
}

func (r *PostgresSegmentRepository) GetSegmentBySlug(slug string) (*models.Segment, error) {
	segment := new(models.Segment)
	err := r.db.QueryRow(selectSegmentBySlug, slug).Scan(&segment.Id, &segment.Slug, &segment.Description)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
	}

	return segment, nil
}

func (r *PostgresSegmentRepository) CheckIfSegmentAlreadyExists(slug string) bool {
	segment := new(models.Segment)
	err := r.db.QueryRow(checkIfSegmentExists, slug).Scan(&segment.Id, &segment.Slug, &segment.Description)

	return err != sql.ErrNoRows
}

func (r *PostgresSegmentRepository) CreateSegment(slug, description string) (string, error) {
	if exists := r.CheckIfSegmentAlreadyExists(slug); !exists {
		row := r.db.QueryRow(createSegment, slug, description)
		if err := row.Scan(&slug); err != nil {
			return "", ErrDatabaseWritingError
		}

		return slug, nil
	} else {
		return "", ErrRecordAlreadyExists
	}
}

func (r *PostgresSegmentRepository) UpdateSegment(slug, description string) (*models.Segment, error) {
	updating := new(models.Segment)
	err := r.db.QueryRow(selectSegmentBySlug, slug).Scan(&updating.Id, &updating.Slug, &updating.Description)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
	}

	_, err = r.db.Exec(updateSegment, description, slug)
	if err != nil {
		return nil, ErrDatabaseWritingError
	}

	updated := &models.Segment{
		Id:          updating.Id,
		Slug:        slug,
		Description: description,
	}

	return updated, nil
}

func (r *PostgresSegmentRepository) DeleteSegment(slug string) (*models.Segment, error) {
	deleted := new(models.Segment)
	err := r.db.QueryRow(selectSegmentBySlug, slug).Scan(&deleted.Id, &deleted.Slug, &deleted.Description)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
	}

	_, err = r.db.Query(deleteSegment, slug)
	if err != nil {
		return nil, ErrDatabaseWritingError
	}

	return deleted, nil
}
