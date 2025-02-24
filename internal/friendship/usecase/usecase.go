package usecase

import (
	"fmt"

	"github.com/lightlink/user-service/internal/friendship/domain/dto"
	friendshipEntity "github.com/lightlink/user-service/internal/friendship/domain/entity"
	friendshipRepository "github.com/lightlink/user-service/internal/friendship/repository"
	groupEntity "github.com/lightlink/user-service/internal/group/domain/entity"
	groupRepository "github.com/lightlink/user-service/internal/group/repository"
	userRepository "github.com/lightlink/user-service/internal/user/repository"
)

type FriendshipUsecaseI interface {
	Create(senderID uint, friendRequest *dto.FriendRequest) (*friendshipEntity.Friendship, error)
	Accept(senderID uint, friendRespond *dto.RespondFriendRequest) (*friendshipEntity.Friendship, error)
	Decline(senderID uint, friendRespond *dto.RespondFriendRequest) (*friendshipEntity.Friendship, error)
	GetPendingRequests(userID uint) ([]*friendshipEntity.Friendship, error)
	GetFriendList(userID uint) ([]*friendshipEntity.Friendship, error)
}

type FriendshipUsecase struct {
	userRepo       userRepository.UserRepositoryI
	friendshipRepo friendshipRepository.FriendshipRepositoryI
	groupRepo      groupRepository.GroupRepositoryI
}

func NewFriendshipUsecase(userRepo userRepository.UserRepositoryI, friendshipRepo friendshipRepository.FriendshipRepositoryI, groupRepo groupRepository.GroupRepositoryI) *FriendshipUsecase {
	return &FriendshipUsecase{
		userRepo:       userRepo,
		friendshipRepo: friendshipRepo,
		groupRepo:      groupRepo,
	}
}

func (uc *FriendshipUsecase) Create(senderID uint, friendRequest *dto.FriendRequest) (*friendshipEntity.Friendship, error) {
	receiverUser, err := uc.userRepo.GetByUsername(friendRequest.ReceiverUseraname)
	if err != nil {
		return nil, err
	}

	user1ID, user2ID := senderID, receiverUser.ID
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	friendshipEntity := &friendshipEntity.Friendship{
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

func (uc *FriendshipUsecase) Accept(senderID uint, friendRespond *dto.RespondFriendRequest) (*friendshipEntity.Friendship, error) {
	fmt.Println("sender_id: ", senderID)
	fmt.Println("receiver_id: ", friendRespond.ReceiverID)
	receiverUser, err := uc.userRepo.GetById(friendRespond.ReceiverID)
	if err != nil {
		return nil, err
	}

	user1ID, user2ID := senderID, receiverUser.ID
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	originalStatus := "accepted"
	compensateStatus := "pending"

	friendship := &friendshipEntity.Friendship{
		User1ID:      user1ID,
		User2ID:      user2ID,
		ActionUserID: senderID,
		StatusName:   originalStatus, /*TODO*/
	}

	_, err = uc.friendshipRepo.Update(friendship)
	if err != nil {
		return nil, err
	}

	err = uc.groupRepo.Create(&groupEntity.PersonalGroup{
		User1ID: user1ID,
		User2ID: user2ID,
	})
	if err != nil {
		friendshipCompensate := &friendshipEntity.Friendship{
			User1ID:      user1ID,
			User2ID:      user2ID,
			ActionUserID: senderID,
			StatusName:   compensateStatus,
		}

		_, compensateErr := uc.friendshipRepo.Update(friendshipCompensate)
		if compensateErr != nil {
			fmt.Println("Critical compensate accepting friend request")
		}

		return nil, err
	}

	return friendship, nil
}

func (uc *FriendshipUsecase) Decline(senderID uint, friendRespond *dto.RespondFriendRequest) (*friendshipEntity.Friendship, error) {
	receiverUser, err := uc.userRepo.GetById(friendRespond.ReceiverID)
	if err != nil {
		return nil, err
	}

	user1ID, user2ID := senderID, receiverUser.ID
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	friendshipEntity := &friendshipEntity.Friendship{
		User1ID:      user1ID,
		User2ID:      user2ID,
		ActionUserID: senderID,
		StatusName:   "declined", /*TODO*/
	}

	_, err = uc.friendshipRepo.Update(friendshipEntity)
	if err != nil {
		return nil, err
	}

	return friendshipEntity, nil
}

func (uc *FriendshipUsecase) GetPendingRequests(userID uint) ([]*friendshipEntity.Friendship, error) {
	pendingRequests, err := uc.friendshipRepo.GetPendingRequests(userID)
	if err != nil {
		return nil, err
	}

	return pendingRequests, nil
}

func (uc *FriendshipUsecase) GetFriendList(userID uint) ([]*friendshipEntity.Friendship, error) {
	friendList, err := uc.friendshipRepo.GetFriendList(userID)
	if err != nil {
		return nil, err
	}

	return friendList, nil
}
