package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/lightlink/user-service/internal/friendship/domain/dto"
	"github.com/lightlink/user-service/internal/friendship/usecase"
)

type FriendshipHandler struct {
	friendshipUC usecase.FriendshipUsecaseI
}

func NewFriendshipHandler(friendshipUsecase usecase.FriendshipUsecaseI) *FriendshipHandler {
	return &FriendshipHandler{
		friendshipUC: friendshipUsecase,
	}
}

func (h *FriendshipHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			/*Handle*/
			fmt.Println(err)
		}
	}()

	friendRequest := &dto.FriendRequest{}
	err = json.Unmarshal(body, friendRequest)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	userIDString := r.Header.Get("X-User-ID")
	userID64, err := strconv.ParseUint(userIDString, 10, 32)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	userID := uint(userID64)

	_, err = h.friendshipUC.Create(userID, friendRequest)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
