package grpc

import (
	"context"

	"github.com/lightlink/user-service/internal/user/domain/dto"
	"github.com/lightlink/user-service/internal/user/domain/entity"
	"github.com/lightlink/user-service/internal/user/usecase"
	proto "github.com/lightlink/user-service/protogen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if err == entity.ErrIsNotExist {
		return nil, status.Error(codes.NotFound, entity.ErrIsNotExist.Error())
	}
	if err != nil {
		return nil, err
	}

	getResponse := dto.EntityToGetUserProto(userEntity)

	return getResponse, nil
}
