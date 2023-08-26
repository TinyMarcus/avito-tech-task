package repositories

import (
	"database/sql"
	"dynamic-user-segmentation-service/internal/db"
	"dynamic-user-segmentation-service/internal/errors"
	"dynamic-user-segmentation-service/internal/models"
	goErrors "errors"
)

type UserRepository interface {
	GetAllUsers() ([]*models.User, error)
	GetUserById(userId string) (*models.User, error)
	CreateUser(name string) error
	DeleteUser(userId string) (*models.User, error)
	AddSegmentToUser(userId, slug, ttl string) error
	TakeSegmentFromUser(userId, slug string) error
	GetActiveSegmentsOfUser(userId string) ([]*models.UsersActiveSegments, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

const (
	selectUsers    = `SELECT id, Name FROM users;`
	selectUserById = `SELECT id, Name FROM users WHERE id = $1;`
	createUser     = `INSERT INTO users (name) VALUES ($1) RETURNING id;`
)

func (r *PostgresUserRepository) GetAllUsers() ([]*models.User, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	var users []*models.User

	rows, err := r.db.Query(selectUsers)
	if err != nil {
		return nil, errors.DatabaseReadingError
	}

	if err := rows.Err(); err != nil {
		return nil, errors.DatabaseReadingError
	}

	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			return nil, errors.DatabaseReadingError
		}
		users = append(users, user)
	}

	defer rows.Close()
	return users, nil
}

func (r *PostgresUserRepository) GetUserById(userId int) (*models.User, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	user := new(models.User)
	err := r.db.QueryRow(selectUserById, userId).Scan(&user.Id, &user.Name)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.RecordNotFound
		}
	}

	return user, nil
}

func (r *PostgresUserRepository) CreateUser(name string) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Query(createUser, name)
	if err != nil {
		return errors.DatabaseWritingError
	}

	return nil
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
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func (r *PostgresUserRepository) AddSegmentToUser(userId int, slug, ttl string) error {
	historyRepository := PostgresHistoryRepository{}
	r.db = db.CreateConnection()
	defer r.db.Close()
	var err error

	user := new(models.User)
	err = r.db.QueryRow(selectUserById, userId).Scan(&user.Id, &user.Name)
	if err != nil {
		return errors.RecordNotFound
	}

	if exists := r.CheckIfUserAlreadyHasSegment(userId, slug); exists == true {
		return nil
	}

	if ttl == "" {
		_, err = r.db.Query(addSegmentToUser, userId, slug, nil)
	} else {
		_, err = r.db.Query(addSegmentToUser, userId, slug, ttl)
	}

	if err != nil {
		return errors.DatabaseWritingError
	}

	historyRepository.SetAddingHistoryRecord(userId, slug)
	return nil
}

func (r *PostgresUserRepository) TakeSegmentFromUser(userId int, slug string) error {
	historyRepository := PostgresHistoryRepository{}
	r.db = db.CreateConnection()
	defer r.db.Close()

	user := new(models.User)
	err := r.db.QueryRow(selectUserById, userId).Scan(&user.Id, &user.Name)
	if err != nil {
		return errors.RecordNotFound
	}

	if exists := r.CheckIfUserAlreadyHasSegment(userId, slug); exists == false {
		return nil
	}

	_, err = r.db.Query(takeSegmentFromUser, userId, slug)
	if err != nil {
		return errors.DatabaseWritingError
	}

	historyRepository.SetRemovingHistoryRecord(userId, slug)
	return nil
}

func (r *PostgresUserRepository) GetActiveSegmentsOfUser(userId int) (*models.UsersActiveSegments, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	user := new(models.User)
	err := r.db.QueryRow(selectUserById, userId).Scan(&user.Id, &user.Name)
	if err != nil {
		return nil, errors.RecordNotFound
	}

	var segments []*models.UserSegment

	rows, err := r.db.Query(getActiveSegmentsOfUser, userId)
	if err != nil {
		return nil, errors.DatabaseReadingError
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		segment := new(models.UserSegment)
		if err := rows.Scan(&segment.UserId, &segment.Slug, &segment.DeadlineDate); err != nil {
			return nil, errors.DatabaseReadingError
		}
		segments = append(segments, segment)
	}

	defer rows.Close()
	return models.ConvertUserSegmentToUsersActiveSegments(userId, segments), nil
}
