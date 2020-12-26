package api

import (
	//"fmt"
	"encoding/json"
	"log"
	"net/http"

	//"github.com/google/uuid"
	"github.com/gorilla/mux"
	mod "github.com/securechat/model"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header()["Date"] = nil
	params := mux.Vars(r)
	roomName := params["roomName"]

	var room mod.ChatRoom
	room.RoomID = roomName

	err := room.CreateChatRoom()
	if err != nil {
		log.Printf("Error Creating room: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header()["Date"] = nil
	params := mux.Vars(r)
	roomName := params["roomName"]

	var room mod.ChatRoom
	room.RoomID = roomName

	result, err := room.DeleteChatRoom()
	if err != nil {
		log.Printf("Error Deleting room: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(result)
}

func GetRoomList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header()["Date"] = nil

	roomList, err := mod.GetAllChatRooms()
	if err != nil {
		log.Printf("Error Getting room list:  %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roomList)
}
