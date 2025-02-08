package dto

import (
	"github.com/lightlink/user-service/internal/user/domain/entity"
	"github.com/lightlink/user-service/internal/user/domain/model"
	proto "github.com/lightlink/user-service/protogen/user"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUserResponse struct {
	ID       uint
	Username string
}

func CreateUserProtoToEntity(createRequest *proto.CreateUserRequest) *entity.User {
	return &entity.User{
		Username: createRequest.Username,
		Password: createRequest.Password,
	}
}

func EntityToGetUserProto(userEntity *entity.User) *proto.GetUserResponse {
	return &proto.GetUserResponse{
		Id:       uint32(userEntity.ID),
		Username: userEntity.Username,
	}
}

func EntityToModel(user *entity.User) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Username:     user.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ModelToEntity(user *model.User) *entity.User {
	return &entity.User{
		ID:       user.ID,
		Username: user.Username,
	}
}
