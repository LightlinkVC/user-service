package postgres

import (
	"database/sql"

	"github.com/lightlink/user-service/internal/user/domain/entity"
	"github.com/lightlink/user-service/internal/user/domain/model"
)

type UserPostgresRepository struct {
	DB *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		DB: db,
	}
}

func (repo *UserPostgresRepository) Create(user *model.User) (*model.User, error) {
	err := repo.DB.
		QueryRow("SELECT 1 FROM users WHERE username = $1", user.Username).Scan(new(int))
	if err == nil {
		return nil, entity.ErrAlreadyCreated
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	var lastID uint
	err = repo.DB.QueryRow("INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", user.Username, user.PasswordHash).Scan(&lastID)
	if err != nil {
		return nil, err
	}

	createdUser := &model.User{
		ID:           uint(lastID),
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}

	return createdUser, nil
}

func (repo *UserPostgresRepository) GetById(id uint) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRow("SELECT id, username, password_hash FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, entity.ErrAlreadyCreated
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
