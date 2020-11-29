package api

import (
	//"fmt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	//"github.com/gorilla/mux"
	mod "github.com/securechat/model"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header()["Date"] = nil

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in CreateRoom  : %v\n", r)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("{\"error\":\"%v\"}", r)))
		}
	}()

	qVals := r.URL.Query()
	val, _ := strconv.Atoi(qVals["val"][0])
	a := 5 / val
	log.Printf(" Value a :%v\n", a)

	var room mod.ChatRoom
	room.RoomID = uuid.New().String()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}
