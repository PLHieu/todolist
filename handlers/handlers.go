package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"todolist/db"

	"github.com/gorilla/mux"
)

// Handler contains the handler and all its dependencies.
type Handler struct {
	ts *db.TodoService
	us *db.UserService
}

// NewHandler initialises a new handler, given dependencies.
func NewHandler(ts *db.TodoService, us *db.UserService) *Handler {
	return &Handler{
		ts: ts,
		us: us,
	}
}

// Index is invoked by HTTP GET /.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & a hardcoded message
	resp := &Response{
		Message: "Welcome to the TodoList service!",
	}
	writeResponse(w, http.StatusOK, resp)
}

// UserUpsert is invoked by HTTP POST /users.
func (h *Handler) UserUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: fmt.Errorf("invalid user body:%v", err).Error(),
		})
		return
	}

	// Initialize a user to unmarshal request body into
	var user db.User
	if err := json.Unmarshal(body, &user); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response{
			Error: fmt.Errorf("invalid user body:%v", err).Error(),
		})
		return
	}

	// Call the repository method corresponding to the operation
	user, err = h.us.Upsert(user)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response{
		User: &user,
	})
}

// ListUserByID is invoked by HTTP GET /users/{id}.
func (h *Handler) ListUserTodoByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	user, todos, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error: err.Error(),
		})
		return
	}

	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response{
		User:  user,
		Todos: todos,
	})
}

// Finish Todo is invoked by POST /todos/{id}
func (h *Handler) FinishTodo(w http.ResponseWriter, r *http.Request) {
	todoID := mux.Vars(r)["id"]
	userID := r.URL.Query().Get("user")
	if err := h.us.Exists(userID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
		return
	}
	_, err := h.ts.MakeTodoDone(todoID, userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error: err.Error(),
		})
		return
	}

	user, todos, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: err.Error(),
		})
		return
	}

	writeResponse(w, http.StatusOK, &Response{
		User:  user,
		Todos: todos,
	})
}

// TodoUpsert is invoked by HTTP POST /todos.
func (h *Handler) TodoUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: fmt.Errorf("invalid todo body:%v", err).Error(),
		})
		return
	}

	// Initialize a todo to unmarshal request body into
	var todo db.Todo
	if err := json.Unmarshal(body, &todo); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response{
			Error: fmt.Errorf("invalid todo body:%v", err).Error(),
		})
		return
	}
	if err := h.us.Exists(todo.OwnerID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
		return
	}

	// Call the repository method corresponding to the operation
	updatedTodo := h.ts.Upsert(todo)
	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response{
		Todos: []db.Todo{updatedTodo},
	})
}

// readRequestBody is a helper method that
// allows to read a request body and return any errors.
func readRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return []byte{}, err
	}
	if err := r.Body.Close(); err != nil {
		return []byte{}, err
	}
	return body, err
}
