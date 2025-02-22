package repository

import (
	"github.com/lightlink/user-service/internal/friendship/domain/entity"
	"github.com/lightlink/user-service/internal/friendship/domain/model"
)

type FriendshipRepositoryI interface {
	Create(friendship *entity.Friendship) (*model.Friendship, error)
}
