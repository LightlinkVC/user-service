package entity

type Friendship struct {
	User1ID      uint   `json:"user1_id"`
	User2ID      uint   `json:"user2_id"`
	StatusName   string `json:"status_name"`
	ActionUserID uint   `json:"action_user_id"`
}
