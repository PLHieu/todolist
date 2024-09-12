package db

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Todo contains all the fields for representing a todo.
type Todo struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`
	OwnerID string `json:"owner_id"`
	Status  string `json:"status"`
}

// TodoService contains all the functionality and dependencies for managing todos.
type TodoService struct {
	DB *gorm.DB
	ns NotificationService
}

// NewTodoService initialises a TodoService given its dependencies.
func NewTodoService(db *gorm.DB, ns NotificationService) *TodoService {
	return &TodoService{
		DB: db,
		ns: ns,
	}
}

// Get returns a given todo or error if none exists.
func (ts *TodoService) Get(id string) (*Todo, error) {
	var t Todo
	if r := ts.DB.Where("id = ?", id).First(&t); r.Error != nil {
		return nil, r.Error
	}

	return &t, nil
}

// Upsert creates or updates a todo.
func (ts *TodoService) Upsert(t Todo) Todo {
	var eb Todo
	if r := ts.DB.Where("id = ?", t.ID).First(&eb); r.Error != nil {
		t.ID = uuid.NewString()
		t.Status = Undone.String()
	}
	ts.DB.Save(&t)
	return t
}

// List returns the list of undone todos.
func (ts *TodoService) List() ([]Todo, error) {
	var items []Todo
	if result := ts.DB.Where("status = ?", Undone.String()).Find(&items); result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

// ListByUser returns the list of todos for a given user.
func (ts *TodoService) ListByUser(userID string) ([]Todo, error) {
	var items []Todo
	if result := ts.DB.Where("owner_id = ?", userID).Find(&items); result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

// MakeTodoDone marks it as done.
func (ts *TodoService) MakeTodoDone(todoID, userID string) (*Todo, error) {
	var t Todo
	if r := ts.DB.Where("id = ?", todoID).First(&t); r.Error != nil {
		return nil, fmt.Errorf("no todo found for id %s:%v", todoID, r.Error)
	}
	if t.Status == Done.String() {
		return nil, fmt.Errorf("todo %s is alreay done", todoID)
	}
	t.OwnerID = userID
	t.Status = Done.String()
	sb := ts.Upsert(t)
	if err := ts.ns.NewNoti(sb); err != nil {
		return nil, err
	}

	return &sb, nil
}
