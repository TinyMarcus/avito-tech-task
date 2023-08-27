package repositories

import (
	"database/sql"
	goErrors "errors"

	"github.com/jmoiron/sqlx"

	"github.com/TinyMarcus/avito-tech-task/internal/handlers/dtos"
	"github.com/TinyMarcus/avito-tech-task/internal/models"
)

type PostgresUserRepository struct {
	db *sqlx.DB
	hr *PostgresHistoryRepository
}

func NewUserRepository(db *sqlx.DB, hr *PostgresHistoryRepository) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
		hr: hr,
	}
}

type HistoryRepository interface {
	SetAddingHistoryRecord(userId int, slug string) error
	SetRemovingHistoryRecord(userId int, slug string) error
	// GetHistoryByDate() ([]*models.Segment, error) TODO: сделать получение истории
}

const (
	selectUsers    = `SELECT id, Name FROM users;`
	selectUserById = `SELECT id, Name FROM users WHERE id = $1;`
	createUser     = `INSERT INTO users (name) VALUES ($1) RETURNING id;`
)

func (r *PostgresUserRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User

	rows, err := r.db.Query(selectUsers)
	if err != nil {
		return nil, ErrDatabaseReadingError
	}

	if err := rows.Err(); err != nil {
		return nil, ErrDatabaseReadingError
	}

	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			return nil, ErrDatabaseReadingError
		}
		users = append(users, user)
	}

	defer rows.Close()
	return users, nil
}

func (r *PostgresUserRepository) GetUserById(userId int) (*models.User, error) {
	user := new(models.User)
	err := r.db.QueryRow(selectUserById, userId).Scan(&user.Id, &user.Name)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
	}

	return user, nil
}

func (r *PostgresUserRepository) CreateUser(name string) (int, error) {
	var id int

	row := r.db.QueryRow(createUser, name)
	if err := row.Scan(&id); err != nil {
		return 0, ErrDatabaseWritingError
	}

	return id, nil
}

const (
	addSegmentToUser        = `INSERT INTO users_segments (user_id, slug, deadline_date) VALUES ($1, $2, $3);`
	takeSegmentFromUser     = `DELETE FROM users_segments WHERE user_id = $1 AND slug = $2;`
	getActiveSegmentsOfUser = `SELECT user_id, slug, deadline_date FROM users_segments 
                                    WHERE user_id = $1 AND (deadline_date IS NULL OR deadline_date > CURRENT_TIMESTAMP);`
	checkIfUserHasSegment = `SELECT user_id, slug, deadline_date FROM users_segments 
                                    WHERE user_id = $1 AND slug = $2`
)

func (r *PostgresUserRepository) CheckIfUserAlreadyHasSegment(userId int, slug string) bool {
	userSegment := new(models.UserSegment)
	err := r.db.QueryRow(checkIfUserHasSegment, userId, slug).Scan(&userSegment.UserId, &userSegment.Slug, &userSegment.DeadlineDate)

	return err != sql.ErrNoRows
}

func (r *PostgresUserRepository) AddSegmentToUser(userId int, slug, ttl string) error {
	var err error

	user := new(models.User)
	err = r.db.QueryRow(selectUserById, userId).Scan(&user.Id, &user.Name)
	if err != nil {
		return ErrRecordNotFound
	}

	if exists := r.CheckIfUserAlreadyHasSegment(userId, slug); exists {
		return nil
	}

	if ttl == "" {
		_, err = r.db.Query(addSegmentToUser, userId, slug, nil)
	} else {
		_, err = r.db.Query(addSegmentToUser, userId, slug, ttl)
	}

	if err != nil {
		return ErrDatabaseWritingError
	}

	err = r.hr.SetAddingHistoryRecord(userId, slug)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) TakeSegmentFromUser(userId int, slug string) error {
	user := new(models.User)
	err := r.db.QueryRow(selectUserById, userId).Scan(&user.Id, &user.Name)
	if err != nil {
		return ErrRecordNotFound
	}

	if exists := r.CheckIfUserAlreadyHasSegment(userId, slug); !exists {
		return nil
	}

	_, err = r.db.Query(takeSegmentFromUser, userId, slug)
	if err != nil {
		return ErrDatabaseWritingError
	}

	err = r.hr.SetRemovingHistoryRecord(userId, slug)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) GetActiveSegmentsOfUser(userId int) (*dtos.UsersActiveSegments, error) {
	user := new(models.User)
	err := r.db.QueryRow(selectUserById, userId).Scan(&user.Id, &user.Name)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	var segments []*models.UserSegment

	rows, err := r.db.Query(getActiveSegmentsOfUser, userId)
	if err != nil {
		return nil, ErrDatabaseReadingError
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		segment := new(models.UserSegment)
		if err := rows.Scan(&segment.UserId, &segment.Slug, &segment.DeadlineDate); err != nil {
			return nil, ErrDatabaseReadingError
		}
		segments = append(segments, segment)
	}

	defer rows.Close()
	return dtos.ConvertUserSegmentToUsersActiveSegments(userId, segments), nil
}
