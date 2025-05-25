package repository

import (
	"github.com/lightlink/user-service/internal/user/domain/entity"
	"github.com/lightlink/user-service/internal/user/domain/model"
)

type UserRepositoryI interface {
	Create(user *entity.User) (*model.User, error)
	GetById(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
}
