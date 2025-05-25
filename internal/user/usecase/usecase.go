package usecase

import (
	"github.com/lightlink/user-service/internal/user/domain/dto"
	"github.com/lightlink/user-service/internal/user/domain/entity"
	"github.com/lightlink/user-service/internal/user/repository"
	proto "github.com/lightlink/user-service/protogen/user"
)

type UserUsecaseI interface {
	Create(createUserRequest *proto.CreateUserRequest) (*entity.User, error)
	GetById(id uint) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
}

type UserUsecase struct {
	userRepo repository.UserRepositoryI
}

func NewUserUsecase(repo repository.UserRepositoryI) *UserUsecase {
	return &UserUsecase{
		userRepo: repo,
	}
}

func (uc *UserUsecase) Create(createUserRequest *proto.CreateUserRequest) (*entity.User, error) {
	userEntity, err := dto.CreateUserProtoToEntity(createUserRequest)
	if err != nil {
		return nil, err
	}

	createdUserModel, err := uc.userRepo.Create(userEntity)
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

func (uc *UserUsecase) GetByUsername(username string) (*entity.User, error) {
	userModel, err := uc.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	userEntity := dto.ModelToEntity(userModel)

	return userEntity, nil
}
