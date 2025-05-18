package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/lightlink/user-service/internal/friendship/domain/dto"
	"github.com/lightlink/user-service/internal/friendship/usecase"
	userEntity "github.com/lightlink/user-service/internal/user/domain/entity"
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
	if err == userEntity.ErrIsNotExist {
		/*Handle*/
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FriendshipHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	friendRespond := &dto.RespondFriendRequest{}
	err = json.Unmarshal(body, friendRespond)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	fmt.Println(friendRespond.ReceiverID)

	userIDString := r.Header.Get("X-User-ID")
	userID64, err := strconv.ParseUint(userIDString, 10, 32)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	userID := uint(userID64)

	_, err = h.friendshipUC.Accept(userID, friendRespond)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FriendshipHandler) DeclineFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	friendRespond := &dto.RespondFriendRequest{}
	err = json.Unmarshal(body, friendRespond)
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

	_, err = h.friendshipUC.Decline(userID, friendRespond)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FriendshipHandler) GetFriendList(w http.ResponseWriter, r *http.Request) {
	userIDString := r.Header.Get("X-User-ID")
	userID64, err := strconv.ParseUint(userIDString, 10, 32)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	userID := uint(userID64)

	friendList, err := h.friendshipUC.GetFriendList(userID)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	response, err := json.Marshal(friendList)
	if err != nil {
		/*Handle*/
		fmt.Println("Failed to marshal friend list response")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		fmt.Println("Failed to write friend list response")
	}
}

func (h *FriendshipHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	userIDString := r.Header.Get("X-User-ID")
	fmt.Printf("Getting pending reqeusts for user with id: %s\n", userIDString)
	userID64, err := strconv.ParseUint(userIDString, 10, 32)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	userID := uint(userID64)

	pendingRequests, err := h.friendshipUC.GetPendingRequests(userID)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	response, err := json.Marshal(pendingRequests)
	if err != nil {
		/*Handle*/
		fmt.Println("Failed to marshal pending requests response")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		fmt.Println("Failed to write pending requests response")
	}
}
