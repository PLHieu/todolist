package db

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User contains all the user fields.
type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	PostCode string `json:"post_code"`
	Country  string `json:"country"`
}

type TodoOperationsService interface {
	ListByUser(userID string) ([]Todo, error)
}

// UserService has all the dependencies required for managing users.
type UserService struct {
	DB *gorm.DB
	ts TodoOperationsService
}

// NewUserService initialises the UserService.
func NewUserService(db *gorm.DB, ts TodoOperationsService) *UserService {
	return &UserService{
		DB: db,
		ts: ts,
	}
}

// Get returns a given user or error if none exists.
func (us *UserService) Get(id string) (*User, []Todo, error) {
	var u User
	if r := us.DB.Where("id = ?", id).First(&u); r.Error != nil {
		return nil, nil, fmt.Errorf("no user found for id %s:%v", id, r.Error)
	}
	todos, err := us.ts.ListByUser(id)
	if err != nil {
		return nil, nil, err
	}

	return &u, todos, nil
}

// Exists returns whether a given user exists and returns an error if none found.
func (us *UserService) Exists(id string) error {
	var u User
	if r := us.DB.Where("id = ?", id).First(&u); r.Error != nil {
		return fmt.Errorf("no user found for id %s:%v", id, r.Error)
	}

	return nil
}

// Upsert creates or updates a new user.
func (us *UserService) Upsert(u User) (User, error) {
	var eu User
	if r := us.DB.Where("id = ?", u.ID).First(&eu); r.Error != nil {
		u.ID = uuid.NewString()
	}
	us.DB.Save(&u)

	return u, nil
}
