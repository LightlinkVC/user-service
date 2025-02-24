package dto

type FriendRequest struct {
	ReceiverUseraname string `json:"username"`
}

type RespondFriendRequest struct {
	ReceiverID uint `json:"receiver_id"`
}
