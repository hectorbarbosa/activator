package models

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        int32  `example:"1"`
	Email     string `validate:"required,email" example:"user1@example.org"`
	Name      string `validate:"required" example:"John"`
	NickName  string `validate:"required" example:"greenman"`
	Activated bool
}

func (s *User) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}
