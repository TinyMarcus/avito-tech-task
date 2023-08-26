package repositories

import (
	"database/sql"
	"dynamic-user-segmentation-service/internal/db"
	"dynamic-user-segmentation-service/internal/errors"
	"dynamic-user-segmentation-service/internal/models"
	goErrors "errors"
)

type SegmentRepository interface {
	GetAllSegments() ([]*models.Segment, error)
	GetSegmentBySlug(slug string) (*models.Segment, error)
	CreateSegment(slug, description string) error
	UpdateSegment(slug, description string) error
	DeleteSegment(slug string) (*models.Segment, error)
}

type PostgresSegmentRepository struct {
	db *sql.DB
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
	r.db = db.CreateConnection()
	defer r.db.Close()

	var segments []*models.Segment

	rows, err := r.db.Query(selectSegments)
	if err != nil {
		return nil, errors.DatabaseReadingError
	}

	if err := rows.Err(); err != nil {
		return nil, errors.DatabaseReadingError
	}

	for rows.Next() {
		segment := new(models.Segment)
		if err := rows.Scan(&segment.Id, &segment.Slug, &segment.Description); err != nil {
			return nil, errors.DatabaseReadingError
		}
		segments = append(segments, segment)
	}

	defer rows.Close()
	return segments, nil
}

func (r *PostgresSegmentRepository) GetSegmentBySlug(slug string) (*models.Segment, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	segment := new(models.Segment)
	err := r.db.QueryRow(selectSegmentBySlug, slug).Scan(&segment.Id, &segment.Slug, &segment.Description)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.RecordNotFound
		}
	}

	return segment, nil
}

func (r *PostgresSegmentRepository) CheckIfSegmentAlreadyExists(slug string) bool {
	segment := new(models.Segment)
	err := r.db.QueryRow(checkIfSegmentExists, slug).Scan(&segment.Id, &segment.Slug, &segment.Description)
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func (r *PostgresSegmentRepository) CreateSegment(slug, description string) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	if exists := r.CheckIfSegmentAlreadyExists(slug); exists == false {
		_, err := r.db.Query(createSegment, slug, description)
		if err != nil {
			return errors.DatabaseWritingError
		}

		return nil
	} else {
		return errors.RecordAlreadyExists
	}
}

func (r *PostgresSegmentRepository) UpdateSegment(slug, description string) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	updated := new(models.Segment)
	err := r.db.QueryRow(selectSegmentBySlug, slug).Scan(&updated.Id, &updated.Slug, &updated.Description)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.RecordNotFound
		}
	}

	_, err = r.db.Exec(updateSegment, description, slug)
	if err != nil {
		return errors.DatabaseWritingError
	}

	return nil
}

func (r *PostgresSegmentRepository) DeleteSegment(slug string) (*models.Segment, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	deleted := new(models.Segment)
	err := r.db.QueryRow(selectSegmentBySlug, slug).Scan(&deleted.Id, &deleted.Slug, &deleted.Description)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.RecordNotFound
		}
	}

	_, err = r.db.Query(deleteSegment, slug)
	if err != nil {
		return nil, errors.DatabaseWritingError
	}

	return deleted, nil
}
