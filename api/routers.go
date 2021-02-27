package api

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Action      string
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc
		if strings.Contains(route.Action, "ValidationRequired") {
			handler = Validator(handler)
		}
		handler = Logger(handler, route.Name)
		handler = Recovery(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v1/",
		Index,
		"SkipValidation",
	},
	/*
		Route{
			"CreateRoom",
			"POST",
			"/v1/chat/room/{roomId}",
			CreateRoom,
			"ValidationRequired",
		},
		Route{
			"DeleteRoom",
			"DELETE",
			"/v1/chat/room/{roomId}",
			DeleteRoom,
			"ValidationRequired",
		},
		Route{
			"GetRoomList",
			"GET",
			"/v1/chat/room",
			GetRoomList,
			"ValidationRequired",
		},
	*/
	Route{
		"ChatLaunch",
		"GET",
		"/v1/chat/html",
		ChatLaunch,
		"SkipValidation",
	},
	Route{
		"ChatBroadCast",
		"GET",
		"/v1/ws/chat/broadcast/{roomId}",
		ChatBroadcast,
		"SkipValidation",
		//"ValidationRequired",
	},
	Route{
		"ChatUnicast",
		"GET",
		"/v1/ws/chat/unicast/{roomId}/{userId}",
		ChatUnicast,
		"SkipValidation",
	},
	/*
		Route{
			"DeleteUser",
			"DELETE",
			"/v1/admin/user/{userId}",
			DeleteUser,
			"ValidationRequired",
		},
		Route{
			"CreateUser",
			"POST",
			"/v1/user",
			CreateUser,
			"ValidationRequired",
		},
	*/
}
