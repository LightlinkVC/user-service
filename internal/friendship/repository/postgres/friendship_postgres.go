package postgres

import (
	"database/sql"

	"github.com/lightlink/user-service/internal/friendship/domain/entity"
	"github.com/lightlink/user-service/internal/friendship/domain/model"
)

type FriendshipRepositoryI interface {
	Create(friendship *entity.Friendship) (*model.Friendship, error)
}

type FriendshipPostgresRepository struct {
	DB *sql.DB
}

func NewFriendshipPostgresRepository(db *sql.DB) *FriendshipPostgresRepository {
	return &FriendshipPostgresRepository{
		DB: db,
	}
}

func (repo *FriendshipPostgresRepository) Create(friendship *entity.Friendship) (*model.Friendship, error) {
	var statusID int

	err := repo.DB.QueryRow(
		"SELECT id FROM friendship_statuses WHERE name = $1",
		friendship.StatusName,
	).Scan(&statusID)
	if err != nil {
		return nil, err
	}

	friendshipModel := model.Friendship{}
	err = repo.DB.QueryRow(
		`INSERT INTO friendships (user1_id, user2_id, status_id, action_user_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, user1_id, user2_id, status_id, action_user_id`,
		friendship.User1ID, friendship.User2ID, statusID, friendship.ActionUserID,
	).Scan(&friendshipModel.ID, &friendshipModel.User1ID, &friendshipModel.User2ID, &friendshipModel.StatusID, &friendshipModel.ActionUserID)
	if err != nil {
		return nil, err
	}

	return &friendshipModel, nil
}
