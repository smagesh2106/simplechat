package api

import (
	//"fmt"
	"encoding/json"
	"log"
	"net/http"

	//"github.com/google/uuid"
	"github.com/gorilla/mux"
	mod "github.com/securechat/model"
	"gopkg.in/validator.v2"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header()["Date"] = nil
	params := mux.Vars(r)
	roomId := params["roomId"]

	var room mod.ChatRoom
	room.RoomID = roomId

	//mongo index model works if the value is > 2
	if err := validator.Validate(room); err != nil {
		var errMsg mod.ErrMsg
		errMsg.Error = "Room ID length should be more then 2"
		log.Println(errMsg.Error)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

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
	roomId := params["roomId"]

	if (roomId == "") || (len(roomId) == 0) {
		log.Println("Invalid Room")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var room mod.ChatRoom
	room.RoomID = roomId

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
