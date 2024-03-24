package model

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	ID		  uint64		 `json:"id"`
	Username  string		 `json:"username"`
	Email     string    	 `json:"email"`
	Password  string	     `json:"-"`
	DoB       time.Time      `json:"dob" gorm:"column:dob"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}
type UserSignUp struct {
	Username string		`json:"username" validate:"required"`
	Email    string    	`json:"email" validate:"required,email"`
	Password string    	`json:"password" validate:"required"`
	DoB      string 	`json:"dob" validate:"required"`
}

type UserSignIn struct {
	Email    string    	`json:"email" validate:"required,email"`
	Password string    	`json:"password" validate:"required"`
}

type UserEditReq struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
}

type UserResponse struct {
	ID       uint64		`json:"id"`
	Email    string		`json:"email"`
	Username string		`json:"username"`
	DoB      string		`json:"dob"`
	UpdatedAt time.Time	`json:"updated_at"`
}

func (u UserSignUp) Validate() error {
	validate := validator.New()

	if err := validate.Struct(u); err != nil {
		return err
	}

	if len(u.Password) < 6 {
		return errors.New("password less than 6 characters")
	}
	dob, err := time.Parse("2006-01-02", u.DoB)
	if err != nil {
		return err
	}

	age := time.Since(dob)
	years := int(age.Hours() / 24 / 365)

	if years < 8 {
		return errors.New("age less than 8 years")
	}	
	return nil

}