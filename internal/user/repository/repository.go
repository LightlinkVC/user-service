package repository

import (
	"github.com/lightlink/user-service/internal/user/domain/model"
)

type UserRepositoryI interface {
	Create(user *model.User) (*model.User, error)
	GetById(id uint) (*model.User, error)
}
