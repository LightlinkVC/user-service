package usecase

import (
	"github.com/lightlink/user-service/internal/friendship/domain/dto"
	"github.com/lightlink/user-service/internal/friendship/domain/entity"
	friendshipRepository "github.com/lightlink/user-service/internal/friendship/repository"
	userRepository "github.com/lightlink/user-service/internal/user/repository"
)

type FriendshipUsecaseI interface {
	Create(senderID uint, friendRequest *dto.FriendRequest) (*entity.Friendship, error)
}

type FriendshipUsecase struct {
	userRepo       userRepository.UserRepositoryI
	friendshipRepo friendshipRepository.FriendshipRepositoryI
}

func NewFriendshipUsecase(userRepo userRepository.UserRepositoryI, friendshipRepo friendshipRepository.FriendshipRepositoryI) *FriendshipUsecase {
	return &FriendshipUsecase{
		userRepo:       userRepo,
		friendshipRepo: friendshipRepo,
	}
}

func (uc *FriendshipUsecase) Create(senderID uint, friendRequest *dto.FriendRequest) (*entity.Friendship, error) {
	receiverUser, err := uc.userRepo.GetByUsername(friendRequest.ReceiverUseraname)
	if err != nil {
		return nil, err
	}

	user1ID, user2ID := senderID, receiverUser.ID
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	friendshipEntity := &entity.Friendship{
		User1ID:      user1ID,
		User2ID:      user2ID,
		ActionUserID: senderID,
		StatusName:   "pending", /*TODO*/
	}

	_, err = uc.friendshipRepo.Create(friendshipEntity)
	if err != nil {
		return nil, err
	}

	return friendshipEntity, nil
}
