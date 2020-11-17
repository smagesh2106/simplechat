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
		handler = Validator(Logger(handler, route.Name))

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
		"/v2/",
		Index,
	},
	/*
		Route{
			"CreateUser",
			strings.ToUpper("Post"),
			"/v2/user",
			CreateUser,
		},

		Route{
			"CreateUsersWithArrayInput",
			strings.ToUpper("Post"),
			"/v2/createWithArray",
			CreateUsersWithArrayInput,
		},

		Route{
			"CreateUsersWithListInput",
			strings.ToUpper("Post"),
			"/v2/createWithList",
			CreateUsersWithListInput,
		},

		Route{
			"DeleteUserById",
			strings.ToUpper("Delete"),
			"/v2/user/{id}",
			DeleteUserById,
		},

		Route{
			"FindByAgeHeader",
			strings.ToUpper("Get"),
			"/v2/findByAgeHeader",
			FindByAgeHeader,
		},

		Route{
			"FindByAgeQuery",
			strings.ToUpper("Get"),
			"/v2/findByAgeQuery",
			FindByAgeQuery,
		},

		Route{
			"GetUserById",
			strings.ToUpper("Get"),
			"/v2/user/{id}",
			GetUserById,
		},

		Route{
			"UpdateUserById",
			strings.ToUpper("Put"),
			"/v2/user/{id}",
			UpdateUserById,
		},
		Route{
			"GetUsers",
			strings.ToUpper("Get"),
			"/v2/users",
			GetUsers,
		},
	*/
}
