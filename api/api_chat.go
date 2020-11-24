package api

import (
	//"fmt"
	"net/http"

	"encoding/json"

	"github.com/google/uuid"
	//"github.com/gorilla/mux"
	mod "github.com/securechat/model"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var room mod.ChatRoom

	room.RoomID = uuid.New().String()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}
