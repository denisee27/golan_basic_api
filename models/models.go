package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name" validate:"required,min=3,max=100"`
	Email     string    `json:"email" validate:"required,email" gorm:"unique"`
	Password  string    `json:"password" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type CreateUser struct {
	Data User `json:"data"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateUser(user *User) error {
	return validate.Struct(user)
}
