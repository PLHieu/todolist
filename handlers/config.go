package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ConfigureServer configures the routes of this server and binds handler functions to them
func ConfigureServer(handler *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/").Handler(http.HandlerFunc(handler.Index))

	router.Methods("POST").Path("/users").Handler(http.HandlerFunc(handler.UserUpsert))

	router.Methods("POST").Path("/todos/new").Handler(http.HandlerFunc(handler.TodoUpsert))
	router.Methods("GET").Path("/todos").Handler(http.HandlerFunc(handler.ListUserTodoByID))
	router.Methods("POST").Path("/todos/finish/{id}").Handler(http.HandlerFunc(handler.FinishTodo))

	return router
}
