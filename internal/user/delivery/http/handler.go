package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lightlink/user-service/internal/user/domain/entity"
)

type UserHanlder struct {
}

func NewUserHanlder() *UserHanlder {
	return &UserHanlder{}
}

func generateUserToken(secret, userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 10).Unix(),
		"channels": []string{
			entity.PersonalChannel(userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (h *UserHanlder) InfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling incoming user-info request")
	userIDString := r.Header.Get("X-User-ID")

	token, err := generateUserToken(os.Getenv("TOKEN_KEY"), userIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"channels": map[string]string{
			"personal": entity.PersonalChannel(userIDString),
		},
	})
}
