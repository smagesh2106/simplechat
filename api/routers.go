package api

import (
	"fmt"
	"net/http"

	//"strings"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		handler = Validator(handler)

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
	},
	Route{
		"CreateRoom",
		"POST",
		"/v1/chat/room/{roomName}",
		CreateRoom,
	},
	Route{
		"DeleteRoom",
		"DELETE",
		"/v1/chat/room/{roomName}",
		DeleteRoom,
	},
	Route{
		"GetRoomList",
		"GET",
		"/v1/chat/room",
		GetRoomList,
	},
}
