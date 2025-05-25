package grpc

import (
	"context"

	"github.com/lightlink/user-service/internal/group/domain/entity"
	proto "github.com/lightlink/user-service/protogen/group"
)

type GroupGrpcRepository struct {
	Client proto.GroupServiceClient
}

func NewGroupGrpcRepository(client *proto.GroupServiceClient) *GroupGrpcRepository {
	return &GroupGrpcRepository{
		Client: *client,
	}
}

func (repo *GroupGrpcRepository) Create(personalGroupEntity *entity.PersonalGroup) error {
	createRequest := &proto.CreatePersonalGroupRequest{
		User1Id: uint32(personalGroupEntity.User1ID),
		User2Id: uint32(personalGroupEntity.User2ID),
	}

	_, err := repo.Client.CreatePersonalGroup(context.Background(), createRequest)
	if err != nil {
		return err
	}

	return nil
}
