package dto

import (
	"github.com/lightlink/user-service/internal/user/domain/entity"
	"github.com/lightlink/user-service/internal/user/domain/model"
	proto "github.com/lightlink/user-service/protogen/user"
)

func CreateUserProtoToEntity(createRequest *proto.CreateUserRequest) (*entity.User, error) {
	return &entity.User{
		Username:     createRequest.Username,
		PasswordHash: createRequest.PasswordHash,
	}, nil
}

func EntityToGetUserProto(userEntity *entity.User) *proto.GetUserResponse {
	return &proto.GetUserResponse{
		Id:       uint32(userEntity.ID),
		Username: userEntity.Username,
	}
}

func ModelToEntity(user *model.User) *entity.User {
	return &entity.User{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}
}
