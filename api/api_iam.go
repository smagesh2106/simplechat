package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	mod "github.com/securechat/model"
	"gopkg.in/validator.v2"
)

/*
 * Delete a User.
 */
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	userId := params["userId"]

	if (userId == "") || (len(userId) == 0) {
		log.Println("Invalid UserId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user mod.User
	user.UserID = userId

	result, err := user.DeleteUser()
	if err != nil {
		log.Printf("Error Deleting user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(result)
}

/*
 * Create a User.
 */
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header()["Date"] = nil

	user := mod.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error user input encoding %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//validate user
	if err := validator.Validate(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	user.CreatedAt = time.Now()
	user.LastLogin = time.Now()
	err = (&user).CreateUser()
	if err != nil {
		var errMsg mod.ErrMsg
		log.Println("userid, email has to be unique")
		errMsg.Error = "userid, email has to be unique"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
