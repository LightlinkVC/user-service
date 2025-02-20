package grpc

import (
	"context"

	"github.com/lightlink/user-service/internal/user/domain/dto"
	"github.com/lightlink/user-service/internal/user/usecase"
	proto "github.com/lightlink/user-service/protogen/user"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	userUC usecase.UserUsecaseI
}

func NewUserService(uc usecase.UserUsecaseI) *UserService {
	return &UserService{
		userUC: uc,
	}
}

func (us *UserService) CreateUser(ctx context.Context, createRequest *proto.CreateUserRequest) (*proto.GetUserResponse, error) {
	createdUserEntity, err := us.userUC.Create(createRequest)
	if err != nil {
		return nil, err
	}

	getResponse := dto.EntityToGetUserProto(createdUserEntity)

	return getResponse, nil
}

func (us *UserService) GetUserById(ctx context.Context, getRequest *proto.GetUserByIdRequest) (*proto.GetUserResponse, error) {
	userEntity, err := us.userUC.GetById(uint(getRequest.Id))
	if err != nil {
		return nil, err
	}

	getResponse := dto.EntityToGetUserProto(userEntity)

	return getResponse, nil
}

func (us *UserService) GetUserByUsername(ctx context.Context, getRequest *proto.GetUserByUsernameRequest) (*proto.GetUserResponse, error) {
	userEntity, err := us.userUC.GetByUsername(getRequest.Username)
	if err != nil {
		return nil, err
	}

	getResponse := dto.EntityToGetUserProto(userEntity)

	return getResponse, nil
}
