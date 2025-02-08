package usecase

import (
	"github.com/lightlink/user-service/internal/user/domain/dto"
	"github.com/lightlink/user-service/internal/user/domain/entity"
	"github.com/lightlink/user-service/internal/user/repository"
)

type UserUsecaseI interface {
	Create(userEntity *entity.User) (*entity.User, error)
	GetById(id uint) (*entity.User, error)
}

type UserUsecase struct {
	userRepo repository.UserRepositoryI
}

func NewUserUsecase(repo repository.UserRepositoryI) *UserUsecase {
	return &UserUsecase{
		userRepo: repo,
	}
}

func (uc *UserUsecase) Create(userEntity *entity.User) (*entity.User, error) {
	userModel, err := dto.EntityToModel(userEntity)
	if err != nil {
		return nil, err
	}

	createdUserModel, err := uc.userRepo.Create(userModel)
	if err != nil {
		return nil, err
	}

	createdUserEntity := dto.ModelToEntity(createdUserModel)

	return createdUserEntity, nil
}

func (uc *UserUsecase) GetById(id uint) (*entity.User, error) {
	userModel, err := uc.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	userEntity := dto.ModelToEntity(userModel)

	return userEntity, nil
}
