package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lightlink/user-service/internal/friendship/domain/entity"
	"github.com/lightlink/user-service/internal/friendship/domain/model"
)

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

func (repo *FriendshipPostgresRepository) Update(friendship *entity.Friendship) (*model.Friendship, error) {
	fmt.Println("Updating: ", *friendship)
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
		`UPDATE friendships
		SET status_id = $1
		WHERE user1_id = $2 AND user2_id = $3
		RETURNING id, user1_id, user2_id, status_id, action_user_id`,
		statusID, friendship.User1ID, friendship.User2ID,
	).Scan(&friendshipModel.ID, &friendshipModel.User1ID, &friendshipModel.User2ID, &friendshipModel.StatusID, &friendshipModel.ActionUserID)
	if err != nil {
		return nil, err
	}

	return &friendshipModel, nil
}

/*TODO Test fix request*/
func (repo *FriendshipPostgresRepository) GetPendingRequests(userID uint) ([]*entity.FriendMeta, error) {
	rows, err := repo.DB.Query(`
		SELECT 
			CASE 
				WHEN f.user1_id = $1 THEN f.user2_id
				ELSE f.user1_id
			END AS friend_id,
			u.username
		FROM friendships f
		JOIN friendship_statuses fs ON f.status_id = fs.id
		JOIN users u ON u.id = CASE 
								WHEN f.user1_id = $1 THEN f.user2_id
								ELSE f.user1_id
							  END
		WHERE 
			fs.name = 'pending'
			AND (f.user1_id = $1 OR f.user2_id = $1)
			AND f.action_user_id != $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pendingRequests := []*entity.FriendMeta{}
	for rows.Next() {
		var req entity.FriendMeta
		err := rows.Scan(
			&req.UserID,
			&req.Username,
		)
		if err != nil {
			return nil, err
		}
		pendingRequests = append(pendingRequests, &req)
	}

	return pendingRequests, nil
}

func (repo *FriendshipPostgresRepository) GetFriendList(userID uint) ([]*entity.FriendMeta, error) {
	rows, err := repo.DB.Query(
		`SELECT 
			CASE 
				WHEN f.user1_id = $1 THEN f.user2_id 
				ELSE f.user1_id 
			END AS friend_id,
			u.username
		FROM friendships f
		JOIN friendship_statuses fs ON f.status_id = fs.id
		JOIN users u ON u.id = CASE 
								WHEN f.user1_id = $1 THEN f.user2_id 
								ELSE f.user1_id 
							END
		WHERE 
			(f.user1_id = $1 OR f.user2_id = $1) 
			AND fs.name = 'accepted';
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("Error GetFriendList reading rows")
		}
	}()

	friendList := []*entity.FriendMeta{}
	for rows.Next() {
		currentFriendRecord := &entity.FriendMeta{}
		err := rows.Scan(
			&currentFriendRecord.UserID,
			&currentFriendRecord.Username,
		)
		if err != nil {
			fmt.Println("Failed to select friend list")
			return nil, err
		}

		friendList = append(friendList, currentFriendRecord)
	}

	return friendList, nil
}
