package postgres

import (
	"database/sql"

	"github.com/lightlink/user-service/internal/user/domain/entity"
	"github.com/lightlink/user-service/internal/user/domain/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserPostgresRepository struct {
	DB *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		DB: db,
	}
}

func (repo *UserPostgresRepository) Create(user *entity.User) (*model.User, error) {
	err := repo.DB.
		QueryRow("SELECT 1 FROM users WHERE username = $1", user.Username).Scan(new(int))
	if err == nil {
		return nil, entity.ErrAlreadyCreated
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	createdUser := model.User{}
	err = repo.DB.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, password_hash",
		user.Username, user.PasswordHash).Scan(&createdUser.ID, &createdUser.Username, &createdUser.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (repo *UserPostgresRepository) GetById(id uint) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRow("SELECT id, username, password_hash FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "can't find such user")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserPostgresRepository) GetByUsername(username string) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, entity.ErrIsNotExist
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
